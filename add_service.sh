#!/bin/sh

# go exist
if ! command -v go; then
  echo "go not exist in PATH"
  exit 1
fi

# kitex exist
if ! command -v kitex; then
  echo "try: go install github.com/cloudwego/kitex/tool/cmd/kitex@latest"
  exit 1
fi

# args $1 exist
if [ "$1" = "" ]; then
  echo "no argument"
  exit
fi

# kitex_gen
echo "INFO: generating kitex_gen client/server ... "
mkdir -p kitex_gen
kitex -module delivery-backend -I idl/ idl/"$1".proto

# service gen
echo "INFO: generating kitex_gen service ... "
mkdir -p service/"$1"
cd service/"$1" &&
  kitex \
    -module "delivery-backend" \
    -service "$1" \
    -use delivery-backend/kitex_gen/ -I ../../idl/ ../../idl/"$1".proto
