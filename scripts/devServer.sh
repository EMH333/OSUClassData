#!/bin/bash

# cd to script directory
cd "$(dirname "$0")" || exit

# go to server root
cd ../cmd/server || exit

go build -race -o ../../build/OSUCD-server main.go subjectHandling.go leaderboardPage.go || exit

# start server
DBUSER=root DBPASS=my-secret-pw DEV=true ./../../build/OSUCD-server #go run main.go
