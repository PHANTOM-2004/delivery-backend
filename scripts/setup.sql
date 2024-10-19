

-- 创建数据库
CREATE DATABASE TJ_SE;

-- 创建一个user
CREATE USER 'scarlet'@'localhost' IDENTIFIED BY '2252707';

-- 赋予权限
GRANT ALL PRIVILEGES ON TJ_SE.* TO 'scarlet'@'localhost';

-- 建立管理员表
create table delivery_admin(
  `id` int auto_increment,
  `created_at` datetime not null,
  `updated_at` datetime not null,
  `deleted_at` datetime,
  `admin_name` varchar(50) not null,
  `account` varchar(100) not null unique,
  `password` varchar(50) not null,
  primary key(id)
) engine=innodb default charset=utf8 comment='管理员账户';
-- 为account创建索引
CREATE INDEX idx_account ON delivery_admin(account);
