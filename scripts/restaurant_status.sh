#!/bin/sh

set -x

# 创建一个管理员
curl "http://localhost:8000/api/v1/admin/create?super_token=LTDZ&admin_name=cjq&password=123451234512345&account=9876543210" --request POST -w "\n"
echo "finish: admin create"

# 管理员登录
curl -c cookies.txt -F "account=9876543210" -F "password=123451234512345" http://localhost:8000/api/v1/admin/login -w "\n"
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
  http://localhost:8000/api/v1/admin/merchant/create --request POST

# 商家登录
curl -c cookies.txt \
  -F "account=merchant_test" \
  -F "password=12345678" \
  http://localhost:8000/api/v1/merchant/login -w "\n"

# 查看登陆状态
curl -b cookies.txt http://localhost:8000/api/v1/merchant/login-status --request GET -w "\n"

# 新建餐厅1
curl -b cookies.txt \
  -F "restaurant_name=顶真帧主" \
  -F "address=礼堂" \
  -F "description=悦刻五代" \
  -F "minimum_delivery_amount=100" \
  http://localhost:8000/api/v1/merchant/restaurant --request POST

# 获得餐厅
curl -b cookies.txt http://localhost:8000/api/v1/merchant/restaurants

# 获得餐厅状态
curl -b cookies.txt http://localhost:8000/api/v1/merchant/restaurant/1/status

# 设置 1
curl -b cookies.txt http://localhost:8000/api/v1/merchant/restaurant/1/status/1 \
  --request PUT

# 获得餐厅状态
curl -b cookies.txt http://localhost:8000/api/v1/merchant/restaurant/1/status

# 设置 0
curl -b cookies.txt http://localhost:8000/api/v1/merchant/restaurant/1/status/0 \
  --request PUT

# 获得餐厅状态
curl -b cookies.txt http://localhost:8000/api/v1/merchant/restaurant/1/status
