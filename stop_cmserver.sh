#!/usr/bin/env bash
set -e

if [ ! -f cmserver.pid ]; then
  echo "cmserver.pid not found, cmserver not running?"
  exit 1
fi
pid=$(cat cmserver.pid)

if ! kill -0 $pid > /dev/null 2>&1; then
  echo "process $pid not running, removing cmserver.pid"
  rm -f cmserver.pid
  exit 1
fi

echo "stopping cmserver, pid=$pid ..."
kill $pid
rm -f cmserver.pid
echo "cmserver stopped"