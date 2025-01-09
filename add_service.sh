#!/bin/sh

# args $1 exist
if [ "$1" = "" ]; then
  echo "no argument"
  exit
fi

add_service_kitex() {

  # kitex exist
  if ! command -v kitex; then
    echo "try: go install github.com/cloudwego/kitex/tool/cmd/kitex@latest"
    exit 1
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
}

add_service_cwgo() {
  # prepare
  mkdir -p rpc_gen/kitex_gen
  mkdir -p "service/$1"

  # client
  cd rpc_gen && cwgo client --type RPC --service "$1" --module delivery-backend/rpc_gen -I ../idl --idl "../idl/$1.proto"

  # server
  cd "../service/$1" && cwgo server --type RPC --service "$1" \
    --module "delivery-backend/service/$1" \
    --pass "-use delivery-backend/rpc_gen/kitex_gen" -I ../../idl --idl "../../idl/$1.proto"
}

add_service_cwgo "$1"
