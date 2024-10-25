# Delivery 后端核心系统

## 项目说明

本项目是送餐系统的后端，核心语言为`golang`。后端负责人是陈艺天。

- 与小程序交互
- 与中台（管理员系统以及商家）交互

## 开发注意事项

### Prerequisite

#### Local Debug

- `mariadb`
- `redis`
- `go > version 1.22`
- `mkcert`
- `go-swagger`

#### Using Docker

- `go-swagger`(optional)
- `mkcert`
- `docker-compose`

### Docker

为了方便环境的一致性，采用`docker`运行服务。

- <del>注意修改`docker-compose.yml`的挂载目录</del>
- <del>注意修改配置文件的`redis`以及`mariadb`对应的`host`</del>

经过配置优化，使用给定的默认配置即可，没有必要再修改配置文件。只需要保证自己的`docker`已经登陆，可以通过如下命令进行验证。

```shell
docker login
```


#### 启动全部服务

注意进入`deploy/local`目录执行：

```shell
docker-compose up -d
```

如果该命令由于网络问题失败请见后文网络问题的解决方法，然后再执行此命令。

#### 关闭全部服务

注意进入`deploy/local`目录执行：

```shell
docker-compose down
```

PS:我们对于数据库使用的是持久化卷，再次启动数据库容器并不会丢失数据。
这也导致，当上游更新数据库时，启动数据库并不会重新初始化，造成容器内数据库信息不同步;
所以可以考虑删除卷,也就是关闭服务的时候加上`-v`参数。这样下次启动的时候就是全新的数据卷完成初始化。

```shell
docker-compose down -v
```

当然也可以手动删除卷

```shell
docker volume ls #查看有哪些卷
docker volume rm volumeName #这里替换为对应的卷名称
```

如果说你数据库里边的数据还想保留(但是你不应当这么做，因为你完全可以在下次重新执行你上次的插入方法)所以先不考虑这个的解决方案。

#### 进入容器内部

可以直接从`docker-desktop`进入，像之前`miniob`那样。也可以从终端进入。

- 从终端进入服务端对应的容器:

```shell
docker exec -it test_go_service /bin/bash
```

进入后，注意我们挂载的目录是`/home/Projects`，我们启动后默认工作区正是这个目录。
因此直接运行：

```shell
go run -v . --dockertest  #启动项目
```

如果离开了工作目录, 可以通过如下面命令返回工作目录：

```shell
cd $PROJECT
```

目前来说，处于测试阶段，因此需要开发者手动启动，因为可能存在一些bug导致服务宕机。
请测试者注意查看log信息。

- 从终端进入数据库容器

如果希望查数据库表，也可以从终端进入数据库对应的容器(或者从`docker-desktop`进入):

```shell
docker exec -it test_mariadb /bin/bash
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
# 注意名字，需要与镜像的名字保持一致
docker load -i golang.tar
docker load -i mariadb.tar
docker load -i redis.tar
```

之后操作和前面一样。


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
