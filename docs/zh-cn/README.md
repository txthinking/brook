CLI 是命令行界面, 与之对应的另一个概念是 GUI 图形用户界面. Brook CLI 只是一个独立的文件, 就是说并没有安装的概念, 只需要下载这一个文件到你电脑即可. Brook CLI 文件同时具有服务端和客户端的功能.

你可能需要用 **root** 用户或 **sudo** 来运行下面提到的命令, 如果你不熟悉 Linux, 我们建议你使用 root 用户来进行操作.

## Install via curl

Let's take the v20200909 version downloaded on linux amd64 as an example

```
$ curl -L https://github.com/txthinking/brook/releases/download/v20200909/brook_linux_amd64 -o /usr/bin/brook
$ chmod +x /usr/bin/brook
```

You can get the download link corresponding to your system on the [releases](https://github.com/txthinking/brook/releases) page

## Install via &nbsp; [nami](https://github.com/txthinking/nami)

Install nami

```
$ curl -L https://git.io/getnami | bash && sleep 6 && exec -l $SHELL
```

You may want more information on [nami github page](https://github.com/txthinking/nami)

```
$ nami install github.com/txthinking/brook
```

## Install on Archlinux(maybe outdated)

```
$ pacman -S brook
```
