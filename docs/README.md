CLI(Command line interface), Brook CLI is just one file, in other words, brook has no concept of installation, you only need to download the brook file to your computer. Brook CLI has both server and client functions.

You may need to run the below commands as **root user or sudo**, if you are not very familiar with linux, we recommend you to use root user.

## Install via curl

Let's take the v20200901 version downloaded on linux amd64 as an example

```
$ curl -L https://github.com/txthinking/brook/releases/download/v20200901/brook_linux_amd64 -o /usr/bin/brook
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

## Install on Archlinux

```
$ pacman -S brook
```
