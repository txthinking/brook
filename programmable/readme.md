# How to submit your own scripts to the Brook Script Gallery

Just add an object to [gallery.json](https://github.com/txthinking/brook/blob/master/programmable/gallery.json)

| Key | Type | Description |
| --- | --- | --- |
| name | string | Your script or module name |
| url | string | Your script or module url. It can be placed in the programmable directory of this project, or anywhere else |
| kind | string | one of `dnsserver`/`server`/`module`/`client` |
| ca | bool | Need to install CA or not |
| author | string | Your name |
| author_url | string | Your url |

kind:

- `dnsserver`: script for brook dnsserver, dohserver, dnsserveroverbrook
- `server`: script for brook server, wsserver, wssserver, quicserver
- `module`: module for Brook GUI Client
- `client`: script for ipio and brook.openwrt
