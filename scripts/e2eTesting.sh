#!/bin/bash
# This script is used to run the e2e tests for the project.

# cd to script directory
cd "$(dirname "$0")" || exit

# go up a directory level
cd .. || exit

# start the database
./scripts/startDB.sh || exit

# start the dev server
./scripts/devServer.sh & 

# go into the frontend directory
cd cmd/server/frontend || exit

# build the frontend
node build.js ci|| exit

# run the e2e tests
npm run e2e

# kill the dev server
killall "OSUCD-server"

# go back to the root directory
cd ../../../ || exit

# stop the database
./scripts/stopDB.sh
