#!/bin/sh

set -x

# 新建餐厅1
curl -b cookies.txt \
  -F "restaurant_name=顶真帧主" \
  -F "address=礼堂" \
  -F "description=悦刻五代" \
  -F "minimum_delivery_amount=100" \
  http://localhost:8000/api/v1/merchant/restaurant --request POST

# 新建餐厅2
curl -b cookies.txt \
  -F "restaurant_name=顶真666" \
  -F "address=礼堂" \
  -F "description=悦刻五代" \
  -F "minimum_delivery_amount=100" \
  http://localhost:8000/api/v1/merchant/restaurant --request POST

# 新建餐厅3
curl -b cookies.txt \
  -F "restaurant_name=吉祥馄炖" \
  -F "address=同济大学" \
  -F "description=好的" \
  -F "minimum_delivery_amount=600" \
  http://localhost:8000/api/v1/merchant/restaurant --request POST

# 获得餐厅
curl -b cookies.txt http://localhost:8000/api/v1/merchant/restaurants
