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
go build -o build/OSUCD-server cmd/server/main.go || exit

# Build the frontend
cd cmd/server/frontend || exit
node build.js production || exit

# Copy frontend to build directory
cd ../../.. || exit
cp -r cmd/server/frontend/dist build/frontend || exit

# deploy to server
read -r -p "Deploy to server? [y/N]" response
response=${response,,} # tolower
if [[ $response =~ ^(yes|y) ]] || [[ -z $response ]]; then
    rsync -rzcvg --delete --chown=:OSUCD build/ malloy:/var/www/services/OSUCD || exit
    rsync -rzcv --chmod=F600 sql/ malloy:/tmp/OSUCD-sql || exit
else
    echo "Aborting deploy"
    exit 1
fi

# restart server service
read -r -p "Restart server service? [y/N]" response
response=${response,,} # tolower
if [[ $response =~ ^(yes|y) ]] || [[ -z $response ]]; then
    ssh malloy -e "sudo systemctl restart EMH-OSUCD" || exit
else
    echo "Did not restarting server service"
    exit 0
fi
