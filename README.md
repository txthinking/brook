# Brook

<!--THEME:github-->
<!--G-R3M673HK5V-->

[🇨🇳 中文](README_ZH.md)

[![Build Status](https://travis-ci.org/txthinking/brook.svg?branch=master)](https://travis-ci.org/txthinking/brook)
[![License: GPL v3](https://img.shields.io/badge/License-GPL%20v3-blue.svg)](http://www.gnu.org/licenses/gpl-3.0)

[🤝 Telegram](https://t.me/brookgroup)
[🩸 YouTube](https://www.youtube.com/txthinking)
[❤️ Sponsor](https://github.com/sponsors/txthinking)

A cross-platform network tool designed for developers.

[🗣 Subscribe Announcement](https://t.me/txthinking_news)

<!--TOC-->

## Install

### Install brook command

> [nami](https://github.com/txthinking/nami) can automatically download the command corresponding to your system. If on Windows, run in [Git Bash](https://gitforwindows.org)<br/>
> or<br/>
> If your system is not Linux, MacOS, Windows, or don't want nami, you can download it directly on the [releases](https://github.com/txthinking/brook/releases) page<br/>
> or<br/>
> the one line install script: `bash <(curl https://bash.ooo/brook.sh)`<br/>
> or<br/>
> scripts written by others<br/>
> or<br/>
> Archlinux: `pacman -S brook` (may be outdated)<br/>
> or<br/>
> brew: `brew install brook` (may be outdated)<br/>

Install nami

```
bash <(curl https://bash.ooo/nami.sh)
```

Install brook

```
nami install brook
```

### Install Brook GUI client

-   [iOS & M1 Mac](https://apps.apple.com/us/app/brook-a-cross-platform-proxy/id1216002642)
-   [Android: Brook.apk](https://github.com/txthinking/brook/releases/latest/download/Brook.apk)
-   [Google Play](https://play.google.com/store/apps/details?id=com.soulsinger)
-   [macOS](https://github.com/txthinking/brook/releases/latest/download/Brook.dmg)
-   [Windows](https://github.com/txthinking/brook/releases/latest/download/Brook.exe)
    -   Windows: requires that the latest version of Edge(chromium-based) has been installed<br/>
    -   Windows Security Virus & threat protection: Settings -> Update & Security -> Windows Security -> Virus & threat protection -> Virus & threat protection settings -> manage settings -> Exclusions -> Add or remove exclusions -> Add an exclusion -> File -> Select Brook.exe<br/>
-   [OpenWrt](#gui-for-official-openwrt)
-   Linux: brook cli + [Socks5 Configurator](https://chrome.google.com/webstore/detail/socks5-configurator/hnpgnjkeaobghpjjhaiemlahikgmnghb) or [tun2brook](https://github.com/txthinking/tun2brook)

[How the Brook GUI works](https://www.txthinking.com/talks/articles/brook-en.article)

## brook `subcommand` and `command line arguments`

-   all `subcoommand`: `brook --help`
-   command line arguments of `subommand`: `brook xxx --help`

## brook rule format

There are three types of rule files

-   domain list: One domain name per line, the suffix matches mode. Can be a local file or an HTTPS URL
-   CIDR v4 list: One CIDR per line, which can be a local file or an HTTPS URL
-   CIDR v6 list: One CIDR per line, which can be a local file or an HTTPS URL

Rules file can be used for

-   Server-side: blocking domain name and IP
-   brook dns: bypass, block domain
-   brook tproxy: bypass, block, domain, ip
-   OpenWrt: bypass, block, domain, ip
-   Brook GUI: bypass, block, domain, ip
