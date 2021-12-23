# Brook

## Install via curl

For example, for Linux amd64, run the following commands to install

```
curl -L https://github.com/txthinking/brook/releases/latest/download/brook_linux_amd64 -o /usr/bin/brook
chmod +x /usr/bin/brook
```

> The first command is to use curl to download the brook_linux_amd64 binary file for Linux amd64 to the /usr/bin/.<br/>
> The second command is to use chmod to allow executable permissions.

You can go to the [releases](https://github.com/txthinking/brook/releases) to download and install the latest version of brook for your system

---

## Install via &nbsp; [nami](https://github.com/txthinking/nami)

Install nami

```
curl -L https://raw.githubusercontent.com/txthinking/nami/master/install.sh | bash && sleep 3 && exec -l $SHELL
```

Use nami to install brook, she will automatically download the latest version for your system

```
nami install brook
```

Use nami to install [joker](https://github.com/txthinking/joker), she can run brook as deamon, **optional but recommended**

```
nami install joker
```

Use nami to install [jinbe](https://github.com/txthinking/jinbe), she can add auto-start command at boot, **optional**

```
nami install jinbe
```

---

## Package manager

Archlinux

```
pacman -S brook
```

macOS

```
brew install brook
```
