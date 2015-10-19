package kapacitor

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/influxdb/influxdb/influxql"
	imodels "github.com/influxdb/influxdb/models"
	"github.com/influxdb/kapacitor/expr"
	"github.com/influxdb/kapacitor/models"
	"github.com/influxdb/kapacitor/pipeline"
)

// Number of previous states to remember when computing flapping percentage.
const defaultFlapHistory = 21

// The newest state change is weighted 'weightDiff' times more than oldest state change.
const weightDiff = 1.5

// Maximum weight applied to newest state change.
const maxWeight = 1.2

type AlertHandler func(ad AlertData)

type AlertLevel int

const (
	NoAlert AlertLevel = iota
	InfoAlert
	WarnAlert
	CritAlert
)

func (l AlertLevel) String() string {
	switch l {
	case NoAlert:
		return "noalert"
	case InfoAlert:
		return "INFO"
	case WarnAlert:
		return "WARNING"
	case CritAlert:
		return "CRITICAL"
	default:
		panic("unknown AlertLevel")
	}
}

func (l AlertLevel) MarshalText() ([]byte, error) {
	return []byte(l.String()), nil
}

type AlertData struct {
	Level AlertLevel      `json:"level"`
	Data  influxql.Result `json:"data"`
}

type AlertNode struct {
	node
	a        *pipeline.AlertNode
	endpoint string
	handlers []AlertHandler
	levels   []*expr.StatefulExpr
	history  []AlertLevel
	hIdx     int
	flapping bool
}

// Create a new  AlertNode which caches the most recent item and exposes it over the HTTP API.
func newAlertNode(et *ExecutingTask, n *pipeline.AlertNode) (an *AlertNode, err error) {
	an = &AlertNode{
		node: node{Node: n, et: et},
		a:    n,
	}
	an.node.runF = an.runAlert
	// Construct alert handlers
	an.handlers = make([]AlertHandler, 0)
	if n.Post != "" {
		an.handlers = append(an.handlers, an.handlePost)
	}
	if n.From != "" && len(n.ToList) != 0 {
		an.handlers = append(an.handlers, an.handleEmail)
	}
	// Parse level expressions
	an.levels = make([]*expr.StatefulExpr, CritAlert+1)
	if n.Info != "" {
		tree, err := expr.ParseForType(n.Info, expr.ReturnBool)
		if err != nil {
			return nil, err
		}
		an.levels[InfoAlert] = &expr.StatefulExpr{tree, expr.Functions()}
	}
	if n.Warn != "" {
		tree, err := expr.ParseForType(n.Warn, expr.ReturnBool)
		if err != nil {
			return nil, err
		}
		an.levels[WarnAlert] = &expr.StatefulExpr{tree, expr.Functions()}
	}
	if n.Crit != "" {
		tree, err := expr.ParseForType(n.Crit, expr.ReturnBool)
		if err != nil {
			return nil, err
		}
		an.levels[CritAlert] = &expr.StatefulExpr{tree, expr.Functions()}
	}
	// Configure flapping
	if n.UseFlapping {
		history := n.History
		if history == 0 {
			history = defaultFlapHistory
		}
		if history < 2 {
			return nil, errors.New("alert history count must be >= 2")
		}
		an.history = make([]AlertLevel, history)
		if n.FlapLow > 1 || n.FlapHigh > 1 {
			return nil, errors.New("alert flap thresholds are percentages and should be between 0 and 1")
		}
	}
	return
}

func (a *AlertNode) runAlert() error {
	switch a.Wants() {
	case pipeline.StreamEdge:
		for p, ok := a.ins[0].NextPoint(); ok; p, ok = a.ins[0].NextPoint() {
			l := a.determineLevel(p.Fields, p.Tags)
			if a.a.UseFlapping {
				a.updateFlapping(l)
				if a.flapping {
					continue
				}
			}
			if l > NoAlert {
				batch := models.Batch{
					Name:   p.Name,
					Group:  p.Group,
					Tags:   p.Tags,
					Points: []models.TimeFields{{Time: p.Time, Fields: p.Fields}},
				}

				ad := AlertData{
					l,
					a.batchToResult(batch),
				}
				for _, h := range a.handlers {
					h(ad)
				}
			}
		}
	case pipeline.BatchEdge:
		for b, ok := a.ins[0].NextBatch(); ok; b, ok = a.ins[0].NextBatch() {
			triggered := false
			for _, p := range b.Points {
				l := a.determineLevel(p.Fields, b.Tags)
				if l > NoAlert {
					triggered = true
					if a.a.UseFlapping {
						a.updateFlapping(l)
						if a.flapping {
							break
						}
					}
					ad := AlertData{l, a.batchToResult(b)}
					for _, h := range a.handlers {
						h(ad)
					}
					break
				}
			}
			if !triggered && a.a.UseFlapping {
				a.updateFlapping(NoAlert)
			}
		}
	}
	return nil
}

func (a *AlertNode) determineLevel(fields models.Fields, tags map[string]string) (level AlertLevel) {
	for l, se := range a.levels {
		if se == nil {
			continue
		}
		if pass, err := EvalPredicate(se, fields, tags); pass {
			level = AlertLevel(l)
		} else if err != nil {
			a.logger.Println("E! error evaluating expression:", err)
			return
		} else {
			return
		}
	}
	return
}

func (a *AlertNode) batchToResult(b models.Batch) influxql.Result {
	row := models.BatchToRow(b)
	r := influxql.Result{
		Series: imodels.Rows{row},
	}
	return r
}

func (a *AlertNode) updateFlapping(level AlertLevel) {
	a.history[a.hIdx] = level
	a.hIdx = (a.hIdx + 1) % len(a.history)

	l := len(a.history)
	changes := 0.0
	weight := (maxWeight / weightDiff)
	step := (maxWeight - weight) / float64(l-1)
	for i := 1; i < l; i++ {
		// get current index
		c := (i + a.hIdx) % l
		// get previous index
		p := c - 1
		// check for wrap around
		if p < 0 {
			p = l - 1
		}
		if a.history[c] != a.history[p] {
			changes += weight
		}
		weight += step
	}

	p := changes / float64(l-1)

	if a.flapping && p < a.a.FlapLow {
		a.flapping = false
	} else if !a.flapping && p > a.a.FlapHigh {
		a.flapping = true
	}
}

func (a *AlertNode) handlePost(ad AlertData) {
	b, err := json.Marshal(ad)
	if err != nil {
		a.logger.Println("E! failed to marshal alert data json", err)
		return
	}
	buf := bytes.NewBuffer(b)
	_, err = http.Post(a.a.Post, "application/json", buf)
	if err != nil {
		a.logger.Println("E! failed to POST batch", err)
	}
}

func (a *AlertNode) handleEmail(ad AlertData) {
	b, err := json.Marshal(ad)
	if err != nil {
		a.logger.Println("E! failed to marshal alert data json", err)
		return
	}
	if a.et.tm.SMTPService != nil {
		a.et.tm.SMTPService.SendMail(a.a.From, a.a.ToList, a.a.Subject, string(b))
	} else {
		a.logger.Println("W! smtp service not enabled, cannot send email.")
	}
}
