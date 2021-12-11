#!/bin/bash

# cd to script directory
cd "$(dirname "$0")" || exit

# go to server root
cd ../cmd/server || exit

# start server
DBUSER=root DBPASS=my-secret-pw go run main.go
