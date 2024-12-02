#!/bin/sh

set -x

# 创建商家账号
curl -b cookies.txt \
  -F "account=merchant_test" \
  -F "password=12345678" \
  -F "merchant_name=孙智城" \
  -F "phone_number=+8618520192763" \
  -F "merchant_application_id=1" \
  http://localhost:8000/api/v1/admin/merchant/create --request POST

# 商家登录
curl -c cookies.txt -X POST -F "account=merchant_test" \
  -F "password=12345678" \
  http://localhost:8000/api/v1/merchant/login -w "\n"

# 查看登陆状态
curl -b cookies.txt http://localhost:8000/api/v1/merchant/login-status --request GET -w "\n"
