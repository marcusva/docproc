#!/bin/sh
set -e

echo "Creating queues manually to speed up tests..."

curl -s -X POST http://docproc.fileinput:4151/topic/create?topic=input
curl -s -X POST http://docproc.webinput:4151/topic/create?topic=input
curl -s -X POST http://docproc.preproc:4151/topic/create?topic=preprocessed
curl -s -X POST http://docproc.renderer:4151/topic/create?topic=rendered

echo "Starting tests in 10 seconds..."

sleep 10

curl -s -X POST -H "Content-Type: application/json" --data @/test/data/raw.json http://docproc.webinput:80/receive || (echo "Failed sending raw.json" && exit 1)
curl -s -X POST -F "file=@/test/data/testrecords.csv" http://docproc.webinput:80/upload || (echo "Failed sending testrecords.csv" && exit 1)

fincount=`find /test/output -type f | wc -l`
loopcnt=1
maxtries=60
while [ $fincount -lt 5 ]; do
    echo "Waiting for 5 messages to be finished, current count: $fincount, try: $loopcnt / $maxtries"
    sleep 2
    fincount=`find /test/output -type f | wc -l`
    loopcnt=`expr $loopcnt + 1`
    if [ $loopcnt -ge $maxtries ]; then
        echo "Timeout after $maxtries tries..."
        break
    fi
done

echo "Finished all messages, comparing output..."

tar -C /test -xzf test-results.tar.gz
diff -Nur /test/output /test/test-results

exitcode=$?

if [ $exitcode -ne 0 ]; then
    echo "Tests failed!"
else
    echo "Tests successful"
fi
exit $exitcode
