#!/bin/bash

export NETWORK_NAME="go-nats_network"

COMMAND=$1

case $COMMAND in
    "start")
        docker compose -f ./deployments/docker-compose.yaml up
        ;;
    "stop")
        docker compose -f ./deployments/docker-compose.yaml down
        ;;
    *)
        echo "Usage: $0 [start | stop]"
        ;;
esac