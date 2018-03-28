#!/bin/sh
set -e

CIP="docprocci_docproc"
PNAME="docproc-ci"
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
$DOCKER cp examples/data/testrecords.csv $CIP.fileinput_1:/app/data

sleep 10

# DO NOT USE: the following lines are to sync proper results with the test result dir
# $DOCKER exec $CIP.output_1 ls -al /app/output
# $DOCKER cp $CIP.output_1:/app/output/. ./test/results

$DOCKER cp ./test/test-results.tar.gz $CIP.output_1:/app
$DOCKER exec $CIP.output_1 tar -C /app -xzf test-results.tar.gz
$DOCKER exec -it $CIP.output_1 diff -Nur /app/output /app/test-results
exitcode=$?

$DOCKER_COMPOSE -p $PNAME kill
$DOCKER_COMPOSE -p $PNAME rm -f

if [ $exitcode -ne 0 ]; then
    echo "Tests failed!"
else
    echo "Tests successful"
fi
exit $exitcode