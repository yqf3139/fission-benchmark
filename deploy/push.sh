#!/bin/bash

set -e

tag=$1
if [ -z "$tag" ]
then
    tag=latest
fi

. build.sh

docker build -t fission-benchmark-bundle .
docker tag fission-benchmark-bundle yqf3139/fission-benchmark-bundle:$tag
docker push yqf3139/fission-benchmark-bundle:$tag
