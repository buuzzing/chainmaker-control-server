#!/usr/bin/env bash
set -e

echo "starting cmserver ..."
nohup ./cmserver > cmserver.log 2>&1 &
echo $! > cmserver.pid
echo "cmserver started, pid=$(cat cmserver.pid)"