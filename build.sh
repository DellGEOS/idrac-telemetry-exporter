#!/bin/bash

echo "Building executable"
GOOS=linux go build -o build/idrac-telemetry-exporter main.go

echo "Building docker container"
docker build -t idrac-telemetry-exporter .

echo "Pushing new docker container"
docker image tag idrac-telemetry-exporter dbblackdiamond/idrac-telemetry-exporter:v0.1
docker image tag idrac-telemetry-exporter dbblackdiamond/idrac-telemetry-exporter:latest
docker push --all-tags dbblackdiamond/idrac-telemetry-exporter

echo "Removing old container, creating and starting new container"
docker stop idrac-telemetry-exporter
docker rm idrac-telemetry-exporter
docker create --name idrac-telemetry-exporter -p 3355:3355 -v /home/bertrand/containers/idrac-telemetry-exporter/config.yml:/etc/idrac-telemetry-exporter/config.yml dbblackdiamond/idrac-telemetry-exporter:latest
docker start idrac-telemetry-exporter
echo "New container provisioned"
