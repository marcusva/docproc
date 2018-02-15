#!/bin/sh

/bin/nsqd --lookupd-tcp-address=nsqlookupd:4160 > /dev/null 2>&1 &
st=$?
if [ $st -ne 0 ]; then
  echo "Failed to start nsqd"
  exit $st
fi
sleep 3

$APP$ $@ &
st=$?
if [ $st -ne 0 ]; then
  echo "Failed to start $APP$"
  exit $st
fi
while sleep 10; do
  ps aux |grep nsqd |grep -q -v grep
  if [ $? -ne 0 ]; then
    echo "nsqd exited already"
    exit -1
  fi
  ps aux |grep $APP$ |grep -q -v grep
  if [ $? -ne 0 ]; then
    echo "$APP$ exited already"
    exit -1
  fi
done