#!/bin/sh

set -x

# 向餐厅1插入category1
curl -b cookies.txt \
  -F "name=雪豹" \
  -F "type=1" \
  -F "sort=9" \
  -F "status=2" \
  http://localhost:8000/api/v1/merchant/restaurant/1/category \
  --request POST

# 向餐厅1插入category2
curl -b cookies.txt \
  -F "name=桌饺" \
  -F "type=2" \
  -F "sort=1" \
  -F "status=1" \
  http://localhost:8000/api/v1/merchant/restaurant/1/category \
  --request POST

# 获得餐厅1所有category
curl -b cookies.txt \
  http://localhost:8000/api/v1/merchant/restaurant/1/categories

