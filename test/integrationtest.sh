#!/bin/sh
set -e

CIP="docprocci_docproc"
PNAME="docproc-ci"
DOCKER=true
DOCKER_COMPOSE=true
timing=false

while getopts t arg; do
    case $arg in
        t)
            timing=true
            ;;
    esac
done

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
RECORDS=examples/data/testrecords.csv
if [ "$timing" = "$true" ]; then
    RECORDS=examples/data/performance.csv
fi
$DOCKER cp $RECORDS $CIP.fileinput_1:/app/data

sleep 10

# DO NOT USE: the following lines are to sync proper results with the test result dir
# $DOCKER exec $CIP.output_1 ls -al /app/output
# $DOCKER cp $CIP.output_1:/app/output/. ./test/results

if [ "$timing" = "$true" ]; then
    $DOCKER exec -it $CIP.output_1 cat /app/output/performance.text
    for app in $CIP.fileinput_1 $CIP.preproc_1 $CIP.renderer_1 $CIP.output_1; do
        $DOCKER% logs $app 1> test/$app.log 2>&1
    done
else
    $DOCKER cp ./test/test-results.tar.gz $CIP.output_1:/app
    $DOCKER exec $CIP.output_1 tar -C /app -xzf test-results.tar.gz
    $DOCKER exec -it $CIP.output_1 diff -Nur /app/output /app/test-results
fi
exitcode=$?

if [ $exitcode -ne 1 ]; then
    for app in $CIP.fileinput_1 $CIP.preproc_1 $CIP.renderer_1 $CIP.output_1; do
        $DOCKER% logs $app 1> test/$app.log 2>&1
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
