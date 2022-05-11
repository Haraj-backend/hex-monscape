#!/bin/bash

host="localhost:4566"
waiting_in_second=5
max_retry=10

while getopts h:r:s:w: flag; do
    case "${flag}" in
        h) host=${OPTARG};;
        r) max_retry=${OPTARG};;
        s) services=${OPTARG};;
        w) waiting_in_second=${OPTARG};;
    esac
done

if [ -z $services ]; then
    echo "flag -s (waited services, comma separated) are required"
    exit 1
fi

services_ready=false
retry=1
until [ "$services_ready" = "true" ] || [ $retry -gt $max_retry ]; do
    localstack_health=$(curl -s http://$host/health)
    running_svcs=$(echo $localstack_health | jq '.services | to_entries[] | select(.value=="running").key')

    for expected_svc in $(echo $services | tr ',' '\n'); do
        services_ready=false
        for running_svc in $running_svcs; do
            if [ "$running_svc" = "\"$expected_svc\"" ]; then
                services_ready=true
                break
            fi
        done

        if [ "$services_ready" = "false" ]; then
            break
        fi
    done

    if [ "$services_ready" = "false" ]; then
        sleep $waiting_in_second
        ((retry = retry+1))
    fi
done

if [ "$services_ready" = "false" ]; then
    echo "failed"
    exit 1
fi
