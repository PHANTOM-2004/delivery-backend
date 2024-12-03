#! /bin/bash

set -x

# 创建一个管理员, admin_test,123451234512345
curl "http://localhost:8000/api/v1/admin/create?super_token=LTDZ&admin_name=cjq&password=123451234512345&account=admin_test" -X POST

# 管理员登录
curl -c cookies.txt -F "account=admin_test" -F "password=123451234512345" http://localhost:8000/api/v1/admin/login -X POST

# 验证登录状态
curl -b cookies.txt http://localhost:8000/api/v1/admin/login-status

# 利用convert做一个假图片
convert -size 32x32 xc:white empty.jpg

# 利用后门，上传申请表
curl -b cookies.txt -X POST -F "description=这是一条测试用例" \
  -F "email=2755345380@qq.com" \
  -F "phone_number=+8618537775175" \
  -F "license=@empty.jpg" \
  -F "name=陈家庆" \
  http://localhost:8000/api/v1/admin/hack/merchant-application -H "Content-Type: multipart/form-data"

# 获得申请表
curl -b cookies.txt http://localhost:8000/api/v1/admin/merchant-application/1
