# Brook CLI

## 使用 nami 安装 brook

> [nami](https://github.com/txthinking/nami) 她会自动帮你下载适用你系统的最新版 brook 命令文件.<br/>
> 如果你的系统是 Windows, 你需要在 [Git Bash](https://gitforwindows.org) 里面运行<br/>
> 如果你的系统不是 Linux, MacOS, Windows, 你可以去 [releases](https://github.com/txthinking/brook/releases) 自己下载

安装 nami

```
bash <(curl https://bash.ooo/nami.sh)
```

使用 nami 安装 brook

```
nami install brook
```

使用 nami 安装[joker](https://github.com/txthinking/joker), 她可以让 brook 以守护进程运行, 适用于 Unix 系的操作系统, **这是可选的, 但是建议安装**

```
nami install joker
```

使用 nami 安装[jinbe](https://github.com/txthinking/jinbe), 她可以添加开机启动命令, 适用于 Unix 系的操作系统, **这是可选的**

```
nami install jinbe
```

---

## 通过一键脚本安装, 适用于 Linux 系统

```
bash <(curl https://bash.ooo/brook.sh)
```

---

## 通过包管理器

Archlinux

```
pacman -S brook
```

macOS

```
brew install brook
```
