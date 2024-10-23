# Delivery 后端核心系统

## 项目说明

本项目是送餐系统的后端，核心语言为`golang`。后端负责人是陈艺天。

- 与小程序交互
- 与中台（管理员系统以及商家）交互

## 开发注意事项

### Prerequisite

- `mariadb`
- `redis`
- `go > version 1.22`
- `mkcert`
- `go-swagger`

```shell
#注意启动服务
sudo systemctl start mariadb redis
```

### 启动项目

在项目根目录下执行：

```shell
go run -v ./
```

### Docker

为了方便环境的一致性，采用`docker`运行服务。
- 注意修改`docker-compose.yml`的挂载目录
- 注意修改配置文件的`redis`以及`mariadb`对应的`host`

#### 启动全部服务

注意进入`deploy/local`目录执行：

```shell
docker-compose up -d
```

#### 关闭全部服务

注意进入`deploy/local`目录执行：

```shell
docker-compose down
```

#### 进入容器内部

可以直接从`docker-desktop`进入，像之前`miniob`那样。也可以从终端进入。

从终端进入服务端对应的容器:

```shell
docker exec -it test_mariadb /bin/bash
```

进入后，注意我们挂载的目录是`/home/Projects`

```shell
cd /home/Projects && go run -v ./
```

即可启动项目。


如果希望查数据库表，也可以从终端进入数据库对应的容器(或者从`docker-desktop`进入):

```shell
docker exec -it test_go_service /bin/bash
```

在容器内部启动数据库客户端:

```shell
mariadb -u scarlet -p
```

然后输入密码, 注意密码不会有回显。默认密码是`2252707`。之后就是熟悉的`SQL-client`

#### 网络问题

理论上是一条命令解决的事情，但是难免在国内开发会有幽默网络问题。
如果`docker`出现网络问题，来我这里拷贝镜像。然后执行:

```shell
# 代表把这三个镜像加载到你的本地
docker load -i local-go_service.tar
docker load -i mariadb.tar
docker load -i redis.tar
```

之后操作和前面一样。

#### 端口绑定问题

如果在`docker`内部启动项目失败，出现`3306`端口占用。说明你安装过`mysql`,把他禁用即可。
具体怎么禁用可以`Google`一下(很简单)。

### 接口文档

项目根目录执行

```shell
make
```

通过`make`启动swagger服务，可以在浏览器中查看`api`文档。

注意前置`swagger`，实际上在`Makefile`中也有检测与安装的规则。

```shell
go install github.com/go-swagger/go-swagger/cmd/swagger@latest
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
