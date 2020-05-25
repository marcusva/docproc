#!/bin/sh
set -e

echo "Starting tests..."

curl -s -X POST -H "Content-Type: application/json" --data @/test/data/raw.json http://docproc.webinput:80/receive || (echo "Failed sending raw.json" && exit 1)
curl -s -X POST -F "file=@/test/data/testrecords.csv" http://docproc.webinput:80/upload || (echo "Failed sending testrecords.csv" && exit 1)

fincount=`curl -s -X GET http://docproc.renderer:4151/stats?format=json |sed -n  's/.*"finish_count":\(\d\),.*/\1/p'`
loopcnt=1
while [ $fincount -lt 5 ]; do
    echo "Waiting for 5 messages to be finished, current count: $fincount, try: ($loopcnt / 30)"
    sleep 1
    fincount=`curl -s -X GET http://docproc.renderer:4151/stats?format=json |sed -n  's/.*"finish_count":\(\d\),.*/\1/p'`
    loopcnt=`expr $loopcnt + 1`
    if [ $loopcnt -ge 30 ]; then
        echo "Timeout after 30 seconds..."
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
