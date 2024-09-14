#!/bin/bash

# cd to script directory
cd "$(dirname "$0")" || exit

# go up a directory level
cd .. || exit

# make tmp directory for sql data
mkdir -p /tmp/OSUCD-mysql

# create database docker
podman run --name OSUCD-mysql -p3306:3306 -v /tmp/OSUCD-mysql:/var/lib/mysql \
    -v "$(pwd)"/sql:/tmp/sql \
    -v "$(pwd)"/cmd/classParser:/tmp/classData \
    -e MYSQL_ROOT_PASSWORD=my-secret-pw -d docker.io/library/mysql:8
#docker exec -it OSUCD-mysql bash

# wait till database is up or 1 minute has passed
timeout=60
start_time=$(date +%s)
while ! podman exec -it OSUCD-mysql bash -c "mysql -P3306 -u root -pmy-secret-pw -e \"SELECT 1;\"" >/dev/null 2>&1; do
    sleep 1
    if [ $(( $(date +%s) - start_time )) -gt $timeout ]; then
        echo "Database failed to start"
        exit 1
    fi
done

#create sql database named OSUClassData
#docker exec -it OSUCD-mysql bash -c "mysql -P3306 -u root -pmy-secret-pw"
podman exec -it OSUCD-mysql bash -c "mysql -P3306 -u root -pmy-secret-pw -e \"SOURCE /tmp/sql/initialSchema.sql;\""
podman exec -it OSUCD-mysql bash -c "mysql -P3306 -u root -pmy-secret-pw -e \"SOURCE /tmp/sql/testData.sql;\""
