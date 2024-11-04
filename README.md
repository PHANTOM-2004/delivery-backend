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

#### localhost SSL

对于本地的调试（在服务器部署之前），直接访问[localhost:80](https://localhost:80)是不可行的，因为没有对应证书。因此我们使用[mkcert](https://github.com/FiloSottile/mkcert)提供认证服务。

```shell
# install命令只需要执行一次
mkcert -install
# 生成证书
mkcert -key-file ./localhost-key.pem -cert-file ./localhost-cert.pem localhost
```

同时服务端会启动一个监听`localhost`的服务用于排查，这一服务会在部署时关闭。


### Docker

为了方便环境的一致性，采用`docker`运行服务。

- <del>注意修改`docker-compose.yml`的挂载目录</del>
- <del>注意修改配置文件的`redis`以及`mariadb`对应的`host`</del>

经过配置优化，使用给定的默认配置即可，没有必要再修改配置文件。只需要保证自己的`docker`已经登陆，可以通过如下命令进行验证。

```shell
docker login
```

在这一部分会先说明如何正确启动项目，跟踪服务端。目前经过更改，暂且对于数据库不使用持久卷，方便每次更新的`api`测试。后面的子栏目会继续说明命令。

- 启动服务

注意进入`deploy/local`目录执行：

```shell
docker-compose up -d
```


- 追踪日志(可以替换后面的`test_go_service`为你希望追踪的容器名)： 

```shell
docker logs -f test_go_service
```

- 通过脚本插入数据: 

最后的文件名就是使用的对应脚本，存放在`scripts`文件夹中，注意发送域名是`http://localhost:8000`，这是在容器内部发送，注意看清楚内容。至于插入了什么数据，请自行阅读；前端也可以根据需要添加新的`example`欢迎提出`MR` .

```shell
docker exec test_go_service /bin/sh -c scripts/admin_example.sh
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

- **从终端进入`golang`服务端对应的容器:**

```shell
docker exec -it test_go_service /bin/bash
```

进入后，注意我们挂载的目录是`/home/Projects`，我们启动后默认工作区正是这个目录。
因此直接运行：

```shell
make
# 或者也可以make rundt, 显然前者更方便
```

如果希望停止运行, 直接按下`CTRL+C`。

如果离开了工作目录, 可以通过如下面命令返回工作目录：

```shell
cd $PROJECT
```

目前来说，处于测试阶段，因此需要开发者手动启动，因为可能存在一些bug导致服务宕机。
请测试者注意查看log信息。

- **从终端进入`mariadb`数据库容器**

如果希望查数据库表，也可以从终端进入数据库对应的容器(或者从`docker-desktop`进入):

```shell
docker exec -it test_mariadb /bin/bash
```

在容器内部启动数据库客户端:

```shell
mariadb -u scarlet -p
```

然后输入密码, 注意密码不会有回显。默认密码是`2252707`。之后就是熟悉的`SQL-client`

- **从终端进入`node`容器:**

```shell
docker exec -it test_vite_service /bin/sh
```

PS: 这里默认已经启动了项目，所以无需自己再手动启动。默认启动的时候没有`npm install`，如果报错进入容器进行`npm install`即可。

为了保证`docker`中正确反向代理接受请求，请添加配置:

```javascript
// vite.config.ts
export default defineConfig({
  server: {
    host: "0.0.0.0",
    // 其他配置
  },
  ...
})
```

为了确保在 Docker 开发环境中 HMR（热模块替换）正常工作，你需要在 vite.config.ts 文件中配置 server.watch，启用 usePolling。这样可以让 Vite 在 Docker 环境下及时检测到文件变化。以下是具体的配置示例：

```javascript
// vite.config.ts
export default {
  server: {
    watch: {
      usePolling: true, // 启用轮询
    },
  },
};
```

注意，启用 usePolling 可能会影响性能，特别是在文件数量较多时。
如果你在其他环境中运行而不需要此功能，可以保持该选项关闭。这样可以在修改源码后，确保你能立即看到效果。

> PS:对于任何容器的退出，在容器内部输入`exit`

#### 网络问题

理论上是一条命令解决的事情，但是难免在国内开发会有幽默网络问题。
如果`docker`出现网络问题，来我这里拷贝镜像。然后执行:

```shell
# 代表把这三个镜像加载到你的本地
# 注意名字，需要与镜像的名字保持一致
docker load -i golang.tar
docker load -i mariadb.tar
docker load -i redis.tar
...
```

之后操作和前面一样。

**PS: 我不想管网络问题了，自己解决去，现在用到的镜像越来越多了。**

### 接口文档

- [apifox](https://apifox.com/apidoc/shared-7796c67b-1b9f-4919-9c83-957d81103b31)

- swagger

项目根目录执行

```shell
make
```

通过`make`启动swagger服务，可以在浏览器中查看`api`文档。

注意前置`swagger`，实际上在`Makefile`中也有检测与安装的规则。

```shell
go install github.com/go-swagger/go-swagger/cmd/swagger@latest
```


