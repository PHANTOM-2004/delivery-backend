#!/bin/sh

set -x

# 向餐厅1插入dish 1
curl -b cookies.txt \
  -F "name=雪豹炒鸡扒" \
  -F "price=2089" \
  -F "image=@Makefile.png" \
  -F "description=好的" \
  http://localhost:8000/api/v1/merchant/restaurant/1/dish \
  --request POST

# 向餐厅1插入dish 2
curl -b cookies.txt \
  -F "name=西瓜炖土豆" \
  -F "price=1999" \
  -F "image=@Makefile.png" \
  -F "description=你知道我要说什么" \
  http://localhost:8000/api/v1/merchant/restaurant/1/dish \
  --request POST

# 把dish 1加入category 1
curl -b cookies.txt \
  -d "dishes=1" \
  http://localhost:8000/api/v1/merchant/category/1/dishes/add \
  --request POST

# 把dish 1,2加入category 2
curl -b cookies.txt \
  -d "dishes=1" \
  -d "dishes=2" \
  http://localhost:8000/api/v1/merchant/category/2/dishes/add \
  --request POST

# 获得餐厅1所有category
curl -b cookies.txt \
  http://localhost:8000/api/v1/merchant/restaurant/1/categories

# 获得餐厅1所有dishes
curl -b cookies.txt \
  http://localhost:8000/api/v1/merchant/restaurant/1/dish

# 插入flavors 1
curl -b cookies.txt \
  http://localhost:8000/api/v1/merchant/restaurant/1/flavor/蒜香干拌 \
  --request POST

# 插入flavors 2
curl -b cookies.txt \
  http://localhost:8000/api/v1/merchant/restaurant/1/flavor/红油干拌 \
  --request POST

# 插入flavors 3
curl -b cookies.txt \
  http://localhost:8000/api/v1/merchant/restaurant/1/flavor/别干拌了 \
  --request POST

# 获得所有flavors
curl -b cookies.txt \
  http://localhost:8000/api/v1/merchant/restaurant/1/flavors

# 为菜品1加入flavors
curl -b cookies.txt \
  -d "flavors=1" \
  -d "flavors=2" \
  -d "flavors=3" \
  http://localhost:8000/api/v1/merchant/dish/1/flavors/add \
  --request POST

# 获得餐厅1所有dishes
curl -b cookies.txt \
  http://localhost:8000/api/v1/merchant/restaurant/1/dish
