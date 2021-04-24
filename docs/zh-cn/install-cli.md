# Brook

## 用 curl 直接下载 brook 文件

以 linux 64 位系统为例

```
$ curl -L https://github.com/txthinking/brook/releases/latest/download/brook_linux_amd64 -o /usr/bin/brook
$ chmod +x /usr/bin/brook
```

> 第一条命令是用 curl 命令下载 linux 64 位最新版本的 brook_linux_amd64 命令文件, 并重命名保存到/usr/bin/brook 路径下.<br/>
> 第二条命令是用 chmod 命令赋予/usr/bin/brook 文件可执行权限.

如果你的系统不是 linux 64 位系统, 你可以在[releases](https://github.com/txthinking/brook/releases) 页面找到对应你系统的 brook 文件链接

## 使用[nami](https://github.com/txthinking/nami)安装 brook🔥

安装 nami

```
$ source <(curl -L https://git.io/getnami)
```

使用 nami 安装 brook, 她会自动帮你下载适用你系统的最新版 Brook CLI 文件, 并赋予可执行权限

```
$ nami install github.com/txthinking/brook
```

使用 nami 安装[joker](https://github.com/txthinking/joker), 她可以让 brook 以守护进程运行, **这是可选的, 但是建议安装**

```
$ nami install github.com/txthinking/joker
```

使用 nami 安装[jinbe](https://github.com/txthinking/jinbe), 她可以添加开机启动命令, **这是可选的**

```
$ nami install github.com/txthinking/jinbe
```

> 社区有一个rust版本的[brook-community/jinbe](https://github.com/brook-community/jinbe), 但是参数使用方式可能有些不同
