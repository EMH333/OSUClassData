#!/bin/bash

# This script is used to deploy the application to my "Malloy" server and restart the application.
# Need to do the following:
# - Build the server executable
# - Build the frontend
# - Copy the server executable to the server
# - Copy the frontend to the server
# - Restart the server SystemD service

# cd to "root" directory (one above scripts)
cd "$(dirname "$0")"/.. || exit

# Start from a clean build dir
rm -r build
mkdir -p build
mkdir -p build/frontend/dist

# Build the server executable
#go build -o build/OSUCD-server cmd/server/main.go || exit
#-compiler gccgo --gccgoflags "-static"
#-ldflags="-extldflags=-static -w -s"
docker run --rm -v "$PWD":/usr/src/ethohampton.com/OSUClassData -w /usr/src/ethohampton.com/OSUClassData/cmd/server \
    golang:1.17-alpine \
    sh -c "apk add build-base && go build -v -ldflags='-linkmode=external -extldflags=-static -w -s' -o ../../build/OSUCD-server" || exit
sudo chown "$USER":"$USER" build/OSUCD-server

# Build the frontend
cd cmd/server/frontend || exit
node build.js production || exit

# Copy frontend to build directory
cd ../../.. || exit
cp -r cmd/server/frontend/dist build/frontend || exit

# deploy to server
read -r -p "Deploy to server? [Y/n]" response
response=${response,,} # tolower
if [[ $response =~ ^(yes|y) ]] || [[ -z $response ]]; then
    rsync -rzcvg --delete --chown=:OSUCD build/ malloy:/var/www/services/OSUCD || exit
    rsync -rzcv --chmod=F600 sql/ malloy:/tmp/OSUCD-sql || exit
else
    echo "Aborting deploy"
    exit 1
fi

# restart server service
read -r -p "Restart server service? [Y/n]" response
response=${response,,} # tolower
if [[ $response =~ ^(yes|y) ]] || [[ -z $response ]]; then
    ssh malloy "sudo systemctl restart EMH-OSUCD" || exit
else
    echo "Did not restarting server service"
    exit 0
fi
