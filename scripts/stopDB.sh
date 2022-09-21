#!/bin/bash

#remove data
docker exec -it OSUCD-mysql bash -c "kill -11 \$(pidof mysqld) && rm -r /var/lib/mysql/*"

# stop and remove docker container named OSUCD-mysql
docker stop OSUCD-mysql
docker rm OSUCD-mysql

# remove /tmp/OSUCD-sql
#rm -r /tmp/OSUCD-mysql
