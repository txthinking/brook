# Brook

## Install via curl

Let's take the v20210214 version downloaded on linux amd64 as an example

```
$ curl -L https://github.com/txthinking/brook/releases/download/v20210214/brook_linux_amd64 -o /usr/bin/brook
$ chmod +x /usr/bin/brook
```

> The first command is to use the curl command to download the brook_linux_amd64 command file of the Linux 64-bit v20210214 version, and rename it and save it to the /usr/bin/brook path<br/>
> The second command is to use the chmod command to give executable permissions to the /usr/bin/brook file

You can get the download link corresponding to your system on the [releases](https://github.com/txthinking/brook/releases) page

## Install via &nbsp; [nami](https://github.com/txthinking/nami)

Install nami, You may want more information on [nami github page](https://github.com/txthinking/nami)

```
$ curl -L https://git.io/getnami | bash && sleep 6 && exec -l $SHELL
```

Use nami to install brook, she will automatically download the latest version for your system

```
$ nami install github.com/txthinking/brook
```
