#!/bin/bash

cd "$(dirname "$0")" || exit
docker image build --pull -t ethohampton/osucd-static-build:latest .
