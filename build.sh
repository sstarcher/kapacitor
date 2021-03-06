#!/bin/bash
# Run the build utility via Docker

set -e

# Make sure our working dir is the dir of the script
DIR=$(cd $(dirname ${BASH_SOURCE[0]}) && pwd)
cd $DIR


# Build new docker image
docker build -t influxdata/kapacitor-builder $DIR
echo "Running build.py"
# Run docker
docker run --rm \
    -e AWS_ACCESS_KEY_ID="$AWS_ACCESS_KEY_ID" \
    -e AWS_SECRET_ACCESS_KEY="$AWS_SECRET_ACCESS_KEY" \
    -v $DIR:/gopath/src/github.com/influxdata/kapacitor \
    influxdata/kapacitor-builder \
    "$@"
