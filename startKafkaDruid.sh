#!/bin/bash

# Kafka directory
KAFKA_DIR="/usr/local/kafka"
# Druid directory
DRUID_DIR="/home/hammad/apps/druid"

# Start Kafka
echo "Starting Kafka..."
cd "$KAFKA_DIR"
nohup bin/zookeeper-server-start.sh config/zookeeper.properties > /dev/null 2>&1 &
nohup bin/kafka-server-start.sh config/server.properties > /dev/null 2>&1 &
echo "Kafka started."

# Start Druid
echo "Starting Druid..."
cd "$DRUID_DIR"
nohup ./bin/start-micro-quickstart > /dev/null 2>&1 &
echo "Druid started."

echo "Both Kafka and Druid are now running."
