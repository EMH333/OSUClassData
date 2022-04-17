#!/bin/bash

# cd to script directory
cd "$(dirname "$0")" || exit

# go up a directory level
cd .. || exit

# make tmp directory for sql data
mkdir -p /tmp/OSUCD-mysql

# create database docker
docker run --name OSUCD-mysql -p3306:3306 -v /tmp/OSUCD-mysql:/var/lib/mysql \
    -v "$(pwd)"/sql:/tmp/sql \
    -v "$(pwd)"/cmd/classParser:/tmp/classData \
    -e MYSQL_ROOT_PASSWORD=my-secret-pw -d docker.io/library/mysql:8
#docker exec -it OSUCD-mysql bash

# wait till database is up
while ! docker exec -it OSUCD-mysql bash -c "mysql -P3306 -u root -pmy-secret-pw -e \"SELECT 1;\"" >/dev/null 2>&1; do
    sleep 1
done

#create sql database named OSUClassData
#docker exec -it OSUCD-mysql bash -c "mysql -P3306 -u root -pmy-secret-pw"
docker exec -it OSUCD-mysql bash -c "mysql -P3306 -u root -pmy-secret-pw -e \"SOURCE /tmp/sql/initialSchema.sql;\""
docker exec -it OSUCD-mysql bash -c "mysql -P3306 -u root -pmy-secret-pw -e \"SOURCE /tmp/sql/testData.sql;\""
