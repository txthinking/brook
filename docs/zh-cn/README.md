# Brook

CLI 是命令行界面, 与之对应的另一个概念是 GUI 图形用户界面. **Brook CLI 只是一个独立的文件, 就是说并没有安装的概念, 只需要下载这一个文件到你电脑即可**. Brook CLI 文件同时具有服务端和客户端的功能.

你可能需要用 **root** 用户或 **sudo** 来运行下面提到的命令, **如果你不是非常熟悉 Linux, 我们建议你使用 ubuntu 和 root 用户来进行操作.**

## 用curl直接下载brook文件

截止目前最新版是v20210101, 以linux 64位系统为例

```
$ curl -L https://github.com/txthinking/brook/releases/download/v20210101/brook_linux_amd64 -o /usr/bin/brook
$ chmod +x /usr/bin/brook
```

如果你的系统不是linux 64位系统, 你可以在[releases](https://github.com/txthinking/brook/releases) 页面找到对应你系统的brook文件链接

## 使用[nami](https://github.com/txthinking/nami)安装brook

安装nami

```
$ curl -L https://git.io/getnami | bash && sleep 6 && exec -l $SHELL
```

如果你想看nami的更多信息, 可以去[nami github 页面](https://github.com/txthinking/nami)

使用nami安装brook, 他会自动帮你选择适用你系统的最新版brook文件

```
$ nami install github.com/txthinking/brook
```
