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
    -e MYSQL_ROOT_PASSWORD=my-secret-pw -d mysql:8
#docker exec -it OSUCD-mysql bash
sleep 1
#create sql database named OSUClassData
#docker exec -it OSUCD-mysql bash -c "mysql -P3306 -u root -pmy-secret-pw"
docker exec -it OSUCD-mysql bash -c "mysql -P3306 -u root -pmy-secret-pw -e \"SOURCE /tmp/sql/initialSchema.sql;\""
docker exec -it OSUCD-mysql bash -c "mysql -P3306 -u root -pmy-secret-pw -e \"SOURCE /tmp/sql/testData.sql;\""

echo "Please run a second time after about 30 seconds to create the OSUClassData database $(pwd)"
