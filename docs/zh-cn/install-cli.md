# Brook

## 用curl直接下载brook文件

截止目前最新版是v20210214, 以linux 64位系统为例

```
$ curl -L https://github.com/txthinking/brook/releases/download/v20210214/brook_linux_amd64 -o /usr/bin/brook
$ chmod +x /usr/bin/brook
```

> 第一条命令是用curl命令下载linux 64位的v20210214版本的brook_linux_amd64命令文件, 并重命名保存到/usr/bin/brook路径下.<br/>
> 第二条命令是用chmod命令赋予/usr/bin/brook文件可执行权限.

如果你的系统不是linux 64位系统, 你可以在[releases](https://github.com/txthinking/brook/releases) 页面找到对应你系统的brook文件链接

## 使用[nami](https://github.com/txthinking/nami)安装brook🔥

安装nami

```
$ curl -L https://git.io/getnami | bash && sleep 6 && exec -l $SHELL
```

使用nami安装brook, 她会自动帮你下载适用你系统的最新版Brook CLI文件, 并赋予可执行权限

```
$ nami install github.com/txthinking/brook
```

使用nami安装[joker](https://github.com/txthinking/joker), 她可以让brook以守护进程运行, **这是可选的, 但是建议安装**

```
$ nami install github.com/txthinking/joker
```

使用nami安装[boa](https://github.com/brook-community/boa), 她可以添加开机启动命令, **这是可选的**

```
$ nami install github.com/brook-community/boa
```
