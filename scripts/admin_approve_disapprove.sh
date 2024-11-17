set -x
curl "http://localhost:8000/api/v1/admin/create?super_token=LTDZ&admin_name=cjq&password=123451234512345&account=0123456789" \
    --request POST -w "\n"
curl \
  -F "description=熊大爷" \
  -F "email=1351020091@qq.com" \
  -F "phone_number=+8618537775175" \
  -F "license=@Makefile.png" \
  -F "name=熊爷" \
  http://localhost:8000/api/v1/customer/merchant-application -H "Content-Type: multipart/form-data" \
  --request POST -w "\n"

  # 管理员登录
curl -c cookies.txt -F "account=0123456789" -F "password=123451234512345" http://localhost:8000/api/v1/admin/login -w "\n"
echo "finish: admin login"

# 不通过商家账号
curl -b cookies.txt \
 http://localhost:8000/api/v1/admin/merchant-application/1/disapprove --request PUT

 # 获得申请表
curl -b cookies.txt http://localhost:8000/api/v1/admin/merchant-application/1

# 通过商家账号
curl -b cookies.txt \
 http://localhost:8000/api/v1/admin/merchant-application/1/approve --request PUT

  # 获得申请表
curl -b cookies.txt http://localhost:8000/api/v1/admin/merchant-application/1

# 不通过商家账号
curl -b cookies.txt \
 http://localhost:8000/api/v1/admin/merchant-application/1/disapprove --request PUT

   # 获得申请表
curl -b cookies.txt http://localhost:8000/api/v1/admin/merchant-application/1
