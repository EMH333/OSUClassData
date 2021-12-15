#!/bin/bash

# cd to script directory
cd "$(dirname "$0")" || exit

# go to server root
cd ../cmd/server || exit

go build -o ../../build/OSUCD-server main.go subjectHandling.go || exit

# start server
DBUSER=root DBPASS=my-secret-pw ./../../build/OSUCD-server #go run main.go
