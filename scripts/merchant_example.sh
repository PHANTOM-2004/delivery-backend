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

# 新建餐厅2 --失败
curl -b cookies.txt \
  -F "restaurant_name=顶真帧主" \
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

#修改餐厅顶真帧主 id = 1
curl -b cookies.txt \
  -F "restaurant_name=能大爷水饺" \
  http://localhost:8000/api/v1/merchant/restaurant/1 \
  --request PUT

curl -b cookies.txt \
  -F "address=满天星广场" \
  http://localhost:8000/api/v1/merchant/restaurant/1 \
  --request PUT

# 获得餐厅
curl -b cookies.txt http://localhost:8000/api/v1/merchant/restaurants

# 向餐厅插入category1
curl -b cookies.txt \
  -F "name=雪豹" \
  -F "type=1" \
  -F "sort=9" \
  http://localhost:8000/api/v1/merchant/restaurant/1/category \
  --request POST

# 向餐厅插入category2
curl -b cookies.txt \
  -F "name=桌饺" \
  -F "sort=9" \
  http://localhost:8000/api/v1/merchant/restaurant/1/category \
  --request POST

# 获得餐厅所有category
curl -b cookies.txt \
  http://localhost:8000/api/v1/merchant/restaurant/1/categories

# 修改餐厅的category1
curl -b cookies.txt \
  -F "sort=337" \
  http://localhost:8000/api/v1/merchant/category/1 \
  --request PUT

# 获得餐厅所有category
curl -b cookies.txt \
  http://localhost:8000/api/v1/merchant/restaurant/1/categories

# 向category 1插入dish 1
curl -b cookies.txt \
  -F "name=雪豹炒鸡扒" \
  -F "price=2089" \
  -F "image=@Makefile.png" \
  -F "description=好的" \
  http://localhost:8000/api/v1/merchant/category/1/dish \
  --request POST

# 向category 1插入dish 2
curl -b cookies.txt \
  -F "name=西瓜炖土豆" \
  -F "price=1999" \
  -F "image=@Makefile.png" \
  -F "description=你知道我要说什么" \
  http://localhost:8000/api/v1/merchant/category/1/dish \
  --request POST

# 获得餐厅1所有category
curl -b cookies.txt \
  http://localhost:8000/api/v1/merchant/restaurant/1/categories

# 修改dish 1
curl -b cookies.txt \
  -F "name=现在不叫雪豹了" \
  -F "price=19990" \
  http://localhost:8000/api/v1/merchant/dish/1 \
  --request PUT

# 修改dish 1图片
curl -b cookies.txt \
  -F "image=@Makefile.png" \
  http://localhost:8000/api/v1/merchant/dish/1 \
  --request PUT

# 获得餐厅1所有category
curl -b cookies.txt \
  http://localhost:8000/api/v1/merchant/restaurant/1/categories

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

# 获得dish 1 的 flavors
curl -b cookies.txt \
  http://localhost:8000/api/v1/merchant/dish/1/flavors

# 为菜品1加入dish flavors
curl -b cookies.txt \
  -d "flavors=1" \
  -d "flavors=2" \
  -d "flavors=3" \
  http://localhost:8000/api/v1/merchant/dish/1/flavors/add \
  --request POST

# 获得dish 1的 flavors
curl -b cookies.txt \
  http://localhost:8000/api/v1/merchant/dish/1/flavors

# 修改flavor1号
curl -b cookies.txt \
  http://localhost:8000/api/v1/merchant/flavor/1/name/改名厚乳 --request PUT

# 获得dish 1的 flavors, 此时flavor应当改变
curl -b cookies.txt \
  http://localhost:8000/api/v1/merchant/dish/1/flavors

# 去掉flavor 2, 3
curl -b cookies.txt \
  -d "flavors=2" \
  -d "flavors=3" \
  http://localhost:8000/api/v1/merchant/dish/1/flavors/delete \
  --request POST

# 获得dish 1的 flavors, 此时flavor应当改变
curl -b cookies.txt \
  http://localhost:8000/api/v1/merchant/dish/1/flavors

# 为菜品1再次加入2, 3 dish flavors
curl -b cookies.txt \
  -d "flavors=2" \
  -d "flavors=3" \
  http://localhost:8000/api/v1/merchant/dish/1/flavors/add \
  --request POST

# 获得dish 1的 flavors, 此时flavor应当改变
curl -b cookies.txt \
  http://localhost:8000/api/v1/merchant/dish/1/flavors

# 此时删除1号flavor, 这个时候dish 1的flavor应当改变
curl -b cookies.txt \
  http://localhost:8000/api/v1/merchant/flavor/1 --request DELETE

# 获得所有flavors, 应当改变
curl -b cookies.txt \
  http://localhost:8000/api/v1/merchant/restaurant/1/flavors

# 获得dish 1的 flavors, 此时flavor应当改变
curl -b cookies.txt \
  http://localhost:8000/api/v1/merchant/dish/1/flavors

# 直接删除1号dish
curl -b cookies.txt \
  http://localhost:8000/api/v1/merchant/dish/1 \
  --request DELETE

# 1号dish删除之后，再次获得dish 1的 flavors, 注意去数据库中观测flavor
curl -b cookies.txt \
  http://localhost:8000/api/v1/merchant/dish/1/flavors

# 删除1号category, 理论上此时应当没有菜品
curl -b cookies.txt \
  http://localhost:8000/api/v1/merchant/category/1 \
  --request DELETE

# 获得餐厅1所有category
curl -b cookies.txt \
  http://localhost:8000/api/v1/merchant/restaurant/1/categories

# 删除餐厅1号, 此时数据库应当不存在任何东西
curl -b cookies.txt \
  http://localhost:8000/api/v1/merchant/restaurant/1 \
  --request DELETE

# 再次获得餐厅1所有category
curl -b cookies.txt \
  http://localhost:8000/api/v1/merchant/restaurant/1/categories

set +x
