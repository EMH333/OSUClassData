#!/bin/bash

# cd to script directory
cd "$(dirname "$0")" || exit

# go up a directory level
cd .. || exit

# start server
DBUSER=root DBPASS=my-secret-pw go run cmd/server/main.go
