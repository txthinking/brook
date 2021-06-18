# Brook

## 什么是`CLI`和`GUI`

-   CLI(Command-line interface), 即命令行界面

    -   **进入命令行界面**就可以运行命令, **强烈建议你了解下[命令行三部曲: 幻灯片和视频](https://talks.txthinking.com/), 强烈建议, 一劳永逸**
    -   **通常情况下, 大家在 Linux 服务器上, 使用的都是命令行界面**. 当然 Linux 也有桌面也能运行 GUI
    -   当然, macOS 和 Windows 也有命令行界面, 只是大家可能不常用

-   GUI(Graphical user interface), 即图形用户界面

    -   **通常情况下, 大家用的 macOS/Windows/iOS/Android 上的双击/单击打开的应用都是 GUI**

## Brook 的 CLI 文件和 GUI 文件

**一般情况下**大家使用 Brook 会需要**服务端**和**客户端**结合使用, **当然**Brook CLI 还有**很多其他独立功能**

-   Brook CLI 文件

    -   Brook CLI 文件是一个**独立的命令文件**, 可以说没有安装的概念, 只需要下载这个文件到你电脑, **在命令行界面赋予它可执行权限后运行它即可**
    -   Brook CLI 文件同时具有**服务端功能**和**客户端的功能**, Brook CLI 文件还有**很多其他独立功能**
    -   如果你在 Linux 下, 你又不是非常熟悉 Linux, **强烈建议使用 ubuntu 系统和 root 用户**来运行命令, 切换到 root 用户的命令是 `$ sudo su`
    -   **Brook 的每个功能都只需要一条命令即可, 无需配置**

-   Brook GUI 文件

    -   **保持系统版本最新**
    -   Brook GUI 文件是图形客户端, 只具有**客户端功能**

**Brook 文件说明**, 都可以在[Releases](https://github.com/txthinking/brook/releases)页面下载. 當然你也可以在左邊菜單裏看到用`nami`自動選擇 Brook CLI 文件

| 文件名                             | CLI/GUI | 适用系统                         |
| ---------------------------------- | ------- | -------------------------------- |
| brook_linux_amd64                  | CLI     | Linux 64 位                      |
| brook_linux_386                    | CLI     | Linux 32 位                      |
| brook_linux_arm64                  | CLI     | Linux arm64                      |
| brook_linux_arm7                   | CLI     | Linux arm7                       |
| brook_linux_arm6                   | CLI     | Linux arm6                       |
| brook_linux_arm5                   | CLI     | Linux arm5                       |
| brook_linux_mips                   | CLI     | Linux mips                       |
| brook_linux_mipsle                 | CLI     | Linux mipsle                     |
| brook_linux_mips_softfloat         | CLI     | Linux mips softfloat             |
| brook_linux_mipsle_softfloat       | CLI     | Linux mipsle softfloat           |
| brook_linux_mips64                 | CLI     | Linux mips64                     |
| brook_linux_mips64le               | CLI     | Linux mips64le                   |
| brook_linux_mips64_softfloat       | CLI     | Linux mips64 softfloat           |
| brook_linux_mips64le_softfloat     | CLI     | Linux mips64le softfloat         |
| brook_linux_ppc64                  | CLI     | Linux ppc64                      |
| brook_linux_ppc64le                | CLI     | Linux ppc64le                    |
| brook_freebsd_386                  | CLI     | FreeBSD 32 位                    |
| brook_freebsd_amd64                | CLI     | FreeBSD 64 位                    |
| brook_netbsd_386                   | CLI     | NetBSD 32 位                     |
| brook_netbsd_amd64                 | CLI     | NetBSD 64 位                     |
| brook_openbsd_386                  | CLI     | OpenBSD 32 位                    |
| brook_openbsd_amd64                | CLI     | OpenBSD 64 位                    |
| brook_windows_amd64.exe            | CLI     | Windows 64 位                    |
| brook_windows_386.exe              | CLI     | Windows 32 位                    |
| brook_darwin_amd64.exe             | CLI     | macOS 64 位                      |
| Brook.dmg                          | GUI     | macOS 64 位                      |
| Brook.exe                          | GUI     | Windows 64 位                    |
| Brook.apk                          | GUI     | Android                          |
| brook_linux_amd64.ipk              | GUI     | OpenWrt Linux 64 bit             |
| brook_linux_386.ipk                | GUI     | OpenWrt Linux 32 bit             |
| brook_linux_arm64.ipk              | GUI     | OpenWrt Linux arm64              |
| brook_linux_arm7.ipk               | GUI     | OpenWrt Linux arm7               |
| brook_linux_arm6.ipk               | GUI     | OpenWrt Linux arm6               |
| brook_linux_arm5.ipk               | GUI     | OpenWrt Linux arm5               |
| brook_linux_mips.ipk               | GUI     | OpenWrt Linux mips               |
| brook_linux_mipsle.ipk             | GUI     | OpenWrt Linux mipsle             |
| brook_linux_mips_softfloat.ipk     | GUI     | OpenWrt Linux mips softfloat     |
| brook_linux_mipsle_softfloat.ipk   | GUI     | OpenWrt Linux mipsle softfloat   |
| brook_linux_mips64.ipk             | GUI     | OpenWrt Linux mips64             |
| brook_linux_mips64le.ipk           | GUI     | OpenWrt Linux mips64le           |
| brook_linux_mips64_softfloat.ipk   | GUI     | OpenWrt Linux mips64 softfloat   |
| brook_linux_mips64le_softfloat.ipk | GUI     | OpenWrt Linux mips64le softfloat |
| brook_linux_ppc64.ipk              | GUI     | OpenWrt Linux ppc64              |
| brook_linux_ppc64le.ipk            | GUI     | OpenWrt Linux ppc64le            |
