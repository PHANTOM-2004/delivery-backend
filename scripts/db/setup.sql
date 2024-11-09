-- 目前该部分实际上已弃用且不再维护，因为现在创建数据库通过gorm的automigrate, 不再同时维护两份代码;
--
-- 创建数据库
-- CREATE DATABASE TJ_SE;

-- 创建一个user
-- CREATE USER 'scarlet'@'localhost' IDENTIFIED BY '2252707';

-- 赋予权限
-- GRANT ALL PRIVILEGES ON TJ_SE.* TO 'scarlet'@'localhost';

-- 建立管理员表
create table if not exists
delivery_admin (
    `id` int auto_increment,
    `created_at` datetime not null,
    `updated_at` datetime not null,
    `deleted_at` datetime,
    `admin_name` varchar(50) not null,
    `account` varchar(50) not null unique,
    `password` varchar(50) not null,
    primary key (id)
) engine = innodb default charset = utf8 comment = '管理员账户';
-- 为account创建索引
create index idx_account on delivery_admin (account);


-- 建立商家申请表
create table if not exists
delivery_merchant_application (
    `id` int auto_increment,
    `created_at` datetime not null,
    `updated_at` datetime not null,
    `deleted_at` datetime,
    `status` tinyint default 2 not null comment '1 代表不通过审核，2 代表未审核，3代表通过审核',
    `description` varchar(300) not null comment '申请账号时的简述',
    `license` varchar(200) not null comment '存放了营业执照的路径，需要商家上传图片',
    `email` varchar(50) not null comment '商家接收账号和密码的地址',
    `phone_number` varchar(30) not null comment '申请表创建者的联系方式',
    `name` varchar(20) not null comment '申请表创建者姓名',
    primary key (id)
) engine = innodb default charset = utf8 comment = '商家账户申请表';

-- 建立商家表
create table if not exists
delivery_merchant (
    `id` int auto_increment,
    `created_at` datetime not null,
    `updated_at` datetime not null,
    `deleted_at` datetime,
    `merchant_name` varchar(50) not null,
    `phone_number` varchar(30) not null,
    `account` varchar(50) not null unique,
    `password` varchar(50) not null,
    `status` tinyint default 1 not null comment '1 代表账户状态有效, 0 代表账户状态无效',
    `merchant_application_id` int unique comment '每一个商家账号关联唯一一个商家申请表',
    foreign key (
        merchant_application_id
    ) references delivery_merchant_application (
        id
    ) on delete set null,
    primary key (id)
) engine = innodb default charset = utf8 comment = '商家账户';
-- 为account创建索引
create index idx_account on delivery_merchant (account);

-- 建立商铺表
-- 每一个商家可以有多个商铺，一个商铺必须属于一个对应的商家。
-- https://github.com/go-gorm/datatypes#time, 关于time类可以参照这部分
create table if not exists
delivery_restaurant (
    `id` int auto_increment,
    `created_at` datetime not null,
    `updated_at` datetime not null,
    `deleted_at` datetime,
    `store_name` varchar(50) not null comment '商铺营业状态',
    `address` varchar(100) not null comment '商铺地址',
    `open_time` time not null comment '开张时间',
    `close_time` time not null comment '打烊时间',
    `status` tinyint default 0 not null comment '1代表商铺开业，0代表打烊',
    `description` varchar(300) not null comment '商铺简介',
    `minimum_order_amount` int not null comment '最低起送金额, 小于等于0代表无起送金额',
    `rating` float not null comment '商家店铺评分'
    `delivery_merchant_id` int comment '店铺所属于的商家id',
    foreign key (delivery_merchant_id) references delivery_merchant (
        id
    ) on delete cascade,
    primary key (id)
) engine = innodb default charset = utf8 comment = '外卖店铺';


-- 建立菜品类型表
create table if not exists
delivery_category (
    `id` int auto_increment,
    `created_at` datetime not null,
    `updated_at` datetime not null,
    `deleted_at` datetime,
    `name` varchar(20) not null comment '分类名称，比如素菜，荤菜，小吃，套餐等等',
    `type` tinyint not null comment '分类类型，0代表菜品，1代表套餐',
    `sort` int default 0 not null comment '排序字段，可以通过调整排序字段，调整顺序',
    `status` tinyint not null default 0 comment '状态，1表示启用，0表示禁用。默认禁用',
) engine = innodb default charset = utf8 comment = '菜品类型';
