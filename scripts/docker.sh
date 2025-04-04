#!/bin/bash

export NETWORK_NAME="learn-nats_network"

COMMAND=$1
SUBCOMMAND=$2

function single(){
    cmd=$1
    case $cmd in
        "start")
            docker compose -f ./deployments/docker/single/docker-compose.yaml up
            ;;
        "stop")
            docker compose -f ./deployments/docker/single/docker-compose.yaml down
            ;;
        *)
            echo "Usage: $0 single [start | stop]"
            ;;
    esac
}

case $COMMAND in
    "single")
        single $SUBCOMMAND
        ;;
    *)
        echo "Usage: $0 [single]"
        ;;
esac