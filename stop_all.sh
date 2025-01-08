#!/bin/sh

stop() {
  echo "stopping service at port [:$1]"
  pid=$(lsof -i ":$1" | grep -m 1 'main' | awk '{print $2}')
  if [ ! "$pid" = "" ]; then
    echo "finding pid: [$pid]"
    echo "killing pid: [$pid]"
    kill "$pid"
  fi
}

stop 8000
stop 8001
stop 8002
