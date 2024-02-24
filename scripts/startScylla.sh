#!/bin/bash

aio=$(cat /proc/sys/fs/aio-max-nr)

if [$aio != "1048576"]; then
echo "fs.aio-max-nr = 1048576" >> /etc/sysctl.conf
sysctl -p /etc/sysctl.conf
fi

docker run --name scylla-node1 -d scylladb/scylla:5.2.0 --overprovisioned 1 --smp 1
docker run --name scylla-node2 -d scylladb/scylla:5.2.0 --seeds="$(docker inspect --format='{{ .NetworkSettings.IPAddress }}' scylla-node1)" --overprovisioned 1 --smp 1
docker run --name scylla-node3 -d scylladb/scylla:5.2.0 --seeds="$(docker inspect --format='{{ .NetworkSettings.IPAddress }}' scylla-node1)" --overprovisioned 1 --smp 1

watch docker exec -it scylla-node3 nodetool status -g -n 5 'date +%H:%M'