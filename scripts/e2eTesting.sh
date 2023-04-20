#!/bin/bash
# This script is used to run the e2e tests for the project.
# if the DEV environment variable is set, it won't deal with the dev server or stop the database at the end

# cd to script directory
cd "$(dirname "$0")" || exit

# go up a directory level
cd .. || exit

if [ -z ${DEV+x} ]; then
    # start the database
    ./scripts/startDB.sh || exit

    # start the dev server
    ./scripts/devServer.sh & 
fi

# go into the frontend directory
cd cmd/server/frontend || exit

if [ -z ${DEV+x} ]; then
    # build the frontend
    node build.js ci|| exit
fi

# run the e2e tests
npm run e2e

# go back to the root directory
cd ../../../ || exit

if [ -z ${DEV+x} ]; then
    # kill the dev server
    killall "OSUCD-server"

    # stop the database
    ./scripts/stopDB.sh
fi
