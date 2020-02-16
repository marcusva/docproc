#!/bin/sh
set -e

CIP="docprocperf_docproc"
PNAME="docprocperf"
DOCKER=docker
DOCKER_COMPOSE=docker-compose

echo "Building docker environment..."
$DOCKER build -t docproc/base .
$DOCKER_COMPOSE -p $PNAME build

echo "Starting docker environment..."
$DOCKER_COMPOSE -p $PNAME up -d

echo "Creating queues manually to speed up testing..."
$DOCKER exec -d $CIP.fileinput_1 curl -X POST http://127.0.0.1:4151/topic/create?topic=input
$DOCKER exec -d $CIP.preproc_1 curl -X POST http://127.0.0.1:4151/topic/create?topic=preprocessed
$DOCKER exec -d $CIP.renderer_1 curl -X POST http://127.0.0.1:4151/topic/create?topic=rendered
$DOCKER exec -d $CIP.output_1 curl -X POST http://127.0.0.1:4151/topic/create?topic=output

echo "Starting tests..."
$DOCKER cp examples/data/performance.csv $CIP.fileinput_1:/app/data

sleep 10

$DOCKER exec -it $CIP.output_1 cat /app/output/performance.text
for app in $CIP.fileinput_1 $CIP.preproc_1 $CIP.renderer_1 $CIP.output_1; do
    $DOCKER% logs $app 1> test/$app.log 2>&1
done

$DOCKER_COMPOSE -p $PNAME kill
$DOCKER_COMPOSE -p $PNAME rm -f
