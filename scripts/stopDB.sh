#!/bin/bash

#remove data
podman exec -it OSUCD-mysql bash -c "kill -11 \$(pidof mysqld) && rm -r /var/lib/mysql/*"

# stop and remove docker container named OSUCD-mysql
podman stop OSUCD-mysql
podman rm OSUCD-mysql

# remove /tmp/OSUCD-sql
#rm -r /tmp/OSUCD-mysql
