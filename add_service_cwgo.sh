#!/bin/sh

# args $1 exist
if [ "$1" = "" ]; then
  echo "no argument"
  exit
fi

# prepare
mkdir -p rpc_gen/kitex_gen
mkdir -p "service/$1"

# client
cd rpc_gen && cwgo client --type RPC --service "$1" --module delivery-backend/rpc_gen -I ../idl --idl "../idl/$1.proto"

# server
cd "../service/$1" && cwgo server --type RPC --service "$1" \
  --module "delivery-backend/service/$1" \
  --pass "-use delivery-backend/rpc_gen/kitex_gen" -I ../../idl --idl "../../idl/$1.proto"
