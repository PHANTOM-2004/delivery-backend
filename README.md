# Delivery 后端核心系统

## 项目说明

本项目是送餐系统的后端，核心语言为`golang`。后端负责人是陈艺天。

- 与小程序交互
- 与中台（管理员系统以及商家）交互

## 开发注意事项

### Prerequisite

- `mariadb`
- `redis`
- `go > version 1.20`
- `mkcert`
- `go-swagger`

```shell
#注意启动服务
sudo systemctl start mariadb redis
```

### 接口文档

```shell
make
```

通过`make`启动swagger服务，可以在浏览器中查看`api`文档。

注意前置`swagger`，实际上在`Makefile`中也有检测与安装的规则.

```shell
go install github.com/go-swagger/go-swagger/cmd/swagger
```

### localhost SSL

对于本地的调试（在服务器部署之前），直接访问[localhost:80](https://localhost:80)是不可行的，因为没有对应证书。因此我们使用[mkcert](https://github.com/FiloSottile/mkcert)提供认证服务。

```shell
# install命令只需要执行一次
mkcert -install
# 生成证书
mkcert -key-file ./localhost-key.pem -cert-file ./localhost-cert.pem localhost
```

同时服务端会启动一个监听`localhost`的服务用于排查，这一服务会在部署时关闭。
