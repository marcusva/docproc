#!/bin/sh
set -e

CIP="docprocci_docproc"
PNAME="docprocci"
DOCKER=docker
DOCKER_COMPOSE=docker-compose

echo "Building docker environment..."
$DOCKER build -t docproc/base .
$DOCKER_COMPOSE -p $PNAME build

echo "Starting docker environment..."
$DOCKER_COMPOSE -p $PNAME up -d

echo "Creating queues manually to speed up testing..."
$DOCKER_COMPOSE exec -d $CIP.fileinput_1 curl -X POST http://127.0.0.1:4151/topic/create?topic=input
$DOCKER_COMPOSE exec -d $CIP.webinput_1 curl -X POST http://127.0.0.1:4151/topic/create?topic=input
$DOCKER_COMPOSE exec -d $CIP.preproc_1 curl -X POST http://127.0.0.1:4151/topic/create?topic=preprocessed
$DOCKER_COMPOSE exec -d $CIP.renderer_1 curl -X POST http://127.0.0.1:4151/topic/create?topic=rendered

sleep 5

echo "Starting tests..."
$DOCKER_COMPOSE cp examples/data/testrecords.csv $CIP.fileinput_1:/app/data
$DOCKER_COMPOSE cp examples/data/raw.json $CIP.webinput_1:/raw.json
$DOCKER_COMPOSE exec -d $CIP.webinput_1 curl -X POST -H "Content-Type: application/json" \
    --data @/raw.json http:/localhost/receive

sleep 20

$DOCKER_COMPOSE exec $CIP.output_1 ls -al /app/output

# DO NOT USE: the following lines are to sync proper results with the test result dir
# $DOCKER cp $CIP.output_1:/app/output/. ./test/results

$DOCKER_COMPOSE cp ./test/test-results.tar.gz $CIP.output_1:/app
$DOCKER_COMPOSE exec $CIP.output_1 tar -C /app -xzf test-results.tar.gz
$DOCKER_COMPOSE exec -it $CIP.output_1 diff -Nur /app/output /app/test-results
exitcode=$?

if [ $exitcode -ne 0 ]; then
    for app in $CIP.fileinput_1 $CIP.preproc_1 $CIP.renderer_1 $CIP.output_1; do
        $DOCKER_COMPOSE logs $app
    done
fi

$DOCKER_COMPOSE -p $PNAME kill
$DOCKER_COMPOSE -p $PNAME rm -f

if [ $exitcode -ne 0 ]; then
    echo "Tests failed!"
else
    echo "Tests successful"
fi
exit $exitcode
