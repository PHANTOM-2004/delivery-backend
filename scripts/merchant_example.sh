#!/bin/sh

set -x

# 创建一个管理员
curl "http://localhost:8000/admin/create?super_token=LTDZ&admin_name=cjq&password=123451234512345&account=9876543210" --request POST -w "\n"
echo "finish: admin create"

# 管理员登录
curl -c cookies.txt -F "account=9876543210" -F "password=123451234512345" http://localhost:8000/admin/login -w "\n"
echo "finish: admin login"

# 添加application作为第一条
curl \
  -F "description=ok" \
  -F "email=666@qq.com" \
  -F "phone_number=+8618537775175" \
  -F "license=@Makefile.png" \
  -F "name=szc" \
  http://localhost:8000/api/v1/customer/merchant-application -H "Content-Type: multipart/form-data" \
  --request POST -w "\n"
echo "finish: admin auth"

# 创建商家账号
curl -b cookies.txt \
  -F "account=merchant_test" \
  -F "password=12345678" \
  -F "merchant_name=孙智城" \
  -F "phone_number=+8618520192763" \
  -F "merchant_application_id=1" \
  http://localhost:8000/api/v1/admin/jwt/merchant-create --request POST

# 商家登录
curl -c cookies.txt \
  -F "account=merchant_test" \
  -F "password=12345678" \
  http://localhost:8000/merchant/login -w "\n"

# 商家鉴权
curl -b cookies.txt http://localhost:8000/api/v1/merchant/jwt/auth --request GET -w "\n"

# 查看登陆状态
curl -b cookies.txt http://localhost:8000/api/v1/merchant/jwt/login-status --request GET -w "\n"

# 新建餐厅1
curl -b cookies.txt \
  -F "restaurant_name=顶真帧主" \
  -F "address=礼堂" \
  -F "description=悦刻五代" \
  -F "minimum_delivery_amount=100" \
  http://localhost:8000/api/v1/merchant/jwt/restaurant/create --request POST

# 新建餐厅2 --失败
curl -b cookies.txt \
  -F "restaurant_name=顶真帧主" \
  -F "address=礼堂" \
  -F "description=悦刻五代" \
  -F "minimum_delivery_amount=100" \
  http://localhost:8000/api/v1/merchant/jwt/restaurant/create --request POST

# 新建餐厅3
curl -b cookies.txt \
  -F "restaurant_name=吉祥馄炖" \
  -F "address=同济大学" \
  -F "description=好的" \
  -F "minimum_delivery_amount=600" \
  http://localhost:8000/api/v1/merchant/jwt/restaurant/create --request POST

# 获得餐厅
curl -b cookies.txt http://localhost:8000/api/v1/merchant/jwt/restaurants

#修改餐厅顶真帧主 id = 1
curl -b cookies.txt \
  -F "restaurant_name=能大爷水饺" \
  http://localhost:8000/api/v1/merchant/jwt/restaurant/1/update \
  --request PUT

curl -b cookies.txt \
  -F "address=满天星广场" \
  http://localhost:8000/api/v1/merchant/jwt/restaurant/1/update \
  --request PUT

# 获得餐厅
curl -b cookies.txt http://localhost:8000/api/v1/merchant/jwt/restaurants

# 向餐厅插入category1
curl -b cookies.txt \
  -F "name=雪豹" \
  -F "type=1" \
  -F "sort=9" \
  http://localhost:8000/api/v1/merchant/jwt/restaurant/1/category/create \
  --request POST

# 向餐厅插入category2
curl -b cookies.txt \
  -F "name=桌饺" \
  -F "sort=9" \
  http://localhost:8000/api/v1/merchant/jwt/restaurant/1/category/create \
  --request POST

# 获得餐厅所有category
curl -b cookies.txt \
  http://localhost:8000/api/v1/merchant/jwt/restaurant/1/categories

# 修改餐厅的category1
curl -b cookies.txt \
  -F "sort=337" \
  http://localhost:8000/api/v1/merchant/jwt/restaurant/1/category/1/update \
  --request PUT

# 获得餐厅所有category
curl -b cookies.txt \
  http://localhost:8000/api/v1/merchant/jwt/restaurant/1/categories

set +x