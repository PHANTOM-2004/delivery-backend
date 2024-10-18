

-- 创建数据库
CREATE DATABASE TJ_SE;

-- 创建一个user
CREATE USER 'scarlet'@'localhost' IDENTIFIED BY '2252707';

-- 赋予权限
GRANT ALL PRIVILEGES ON TJ_SE.* TO 'scarlet'@'localhost';

