#!/bin/sh

set -x

# 创建一个管理员
curl "http://localhost:8000/api/v1/admin/create?super_token=LTDZ&admin_name=cjq&password=123451234512345&account=0123456789" --request POST -w "\n"

# 管理员登录
curl -c cookies.txt -F "account=0123456789" -F "password=123451234512345" http://localhost:8000/api/v1/admin/login -w "\n"

# 验证登录状态
curl -b cookies.txt http://localhost:8000/api/v1/admin/login-status

# 管理员修改密码
curl -b cookies.txt http://localhost:8000/api/v1/admin/change-password \
  -F "password=012345012345012345" --request PUT -w "\n"

# 验证登录状态
curl -b cookies.txt http://localhost:8000/api/v1/admin/login-status

# 再次登录
curl -c cookies.txt http://localhost:8000/api/v1/admin/login \
  -F "account=0123456789" \
  -F "password=012345012345012345" --request POST -w "\n"

# 验证登录状态
curl -b cookies.txt http://localhost:8000/api/v1/admin/login-status

# 登出
curl -b cookies.txt http://localhost:8000/api/v1/admin/logout --request POST

# 验证登录状态
curl -b cookies.txt http://localhost:8000/api/v1/admin/login-status

#管理员登录
curl -c cookies.txt -F "account=0123456789" -F "password=012345012345012345" http://localhost:8000/api/v1/admin/login -w "\n"
echo "finish: admin login"

# 验证登录状态
curl -b cookies.txt http://localhost:8000/api/v1/admin/login-status

# 上传申请表
curl \
  -F "description=ok" \
  -F "email=666@qq.com" \
  -F "phone_number=+8618537775175" \
  -F "license=@Makefile.png" \
  -F "name=szc" \
  http://localhost:8000/api/v1/customer/merchant-application -H "Content-Type: multipart/form-data" \
  --request POST -w "\n"
echo "finish: admin auth"

# 获得申请表
curl -b cookies.txt http://localhost:8000/api/v1/admin/merchant-application/1 -w "\n"
echo "finish: admin get merchant-application"

# 批准申请表
curl -b cookies.txt http://localhost:8000/api/v1/admin/merchant-application/1/approve -w "\n" --request PUT

# 获取证书
curl -b cookies.txt http://localhost:8000/api/v1/admin/merchant-application/license/*filepath -w "\n"
echo "finish: admin get merchant license"

# 通过商家申请
# curl -b cookies.txt http://localhost:8000/api/v1/admin/jwt/merchant-application/1/approve --request POST
#
#

set +x
