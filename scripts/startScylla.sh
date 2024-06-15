#!/bin/bash

# Set the desired cluster name
CLUSTER_NAME="ifttt_cluster"

# docker run --name scylla-si -d scylladb/scylla:5.2.0 --smp 2 --memory 4G

check_container() {
  local container_name="$1"
  docker ps -a --filter "name=${container_name}" --format "{{.Names}}"
}

create_or_start_container() {
  local container_name="$1"
  local image_version="5.2.0"

  existing_container=$(check_container "${container_name}")

  if [ -z "${existing_container}" ]; then
    echo "Creating new container: ${container_name}"
    docker run --name "${container_name}" -d \
      -e SCYLLA_CLUSTER=$CLUSTER_NAME \
      scylladb/scylla:${image_version} --overprovisioned 1 --smp 1
  else
    echo "Using existing container: ${existing_container}"
    docker start "${existing_container}"
  fi
}

aio=$(cat /proc/sys/fs/aio-max-nr)

if [ "$aio" != "1048576" ]; then
  echo "fs.aio-max-nr = 1048576" >> /etc/sysctl.conf
  sysctl -p /etc/sysctl.conf
fi

create_or_start_container "scylla-node1"
create_or_start_container "scylla-node2"
create_or_start_container "scylla-node3"

# Wait for all nodes to be up
while true; do
  if docker exec -it scylla-node1 nodetool status | grep "UN" && \
     docker exec -it scylla-node2 nodetool status | grep "UN" && \
     docker exec -it scylla-node3 nodetool status | grep "UN"; then
    echo "All nodes online. Starting cqlsh session"
    sleep 5
    break
  else
    echo "All nodes not online yet"
    sleep 5
  fi
done

docker exec -it scylla-node3 cqlsh