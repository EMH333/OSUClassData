#!/bin/bash

cd "$(dirname "$0")" || exit

# temporarily copy go.mod to build image
cp ../../../go.mod ./ || exit
cp ../../../go.sum ./ || exit

docker image build --pull -t ethohampton/osucd-static-build:latest .

# remove temporary copies
rm go.mod go.sum || exit
