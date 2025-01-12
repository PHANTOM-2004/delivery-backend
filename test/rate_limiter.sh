#!/bin/sh

for i in $(seq 100); do
  curl "http://127.0.0.1:8000/ping"
  sleep 0.3
done
