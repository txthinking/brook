# Brook

## 用curl直接下载brook文件

截止目前最新版是v20210214, 以linux 64位系统为例

```
$ curl -L https://github.com/txthinking/brook/releases/download/v20210214/brook_linux_amd64 -o /usr/bin/brook
$ chmod +x /usr/bin/brook
```

> 第一条命令是用curl命令下载linux 64位的v20210214版本的brook_linux_amd64命令文件, 并重命名保存到/usr/bin/brook路径下<br/>
> 第二条命令是用chmod命令赋予/usr/bin/brook文件可执行权限

如果你的系统不是linux 64位系统, 你可以在[releases](https://github.com/txthinking/brook/releases) 页面找到对应你系统的brook文件链接

## 使用[nami](https://github.com/txthinking/nami)安装brook

安装nami, 如果你想看nami的更多信息, 可以去[nami github 页面](https://github.com/txthinking/nami)

```
$ curl -L https://git.io/getnami | bash && sleep 6 && exec -l $SHELL
```

使用nami安装brook, 他会自动帮你下载适用你系统的最新版Brook CLI文件, 并赋予可执行权限

```
$ nami install github.com/txthinking/brook
```
