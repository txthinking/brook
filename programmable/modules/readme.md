## 模块原理 | Module principle

很简单，一个模块里有几个处理函数。一个头 + N 个模块 + 一个尾，最后合并成一个脚本文件。

Very simple, there are several processing functions in a module. A header + N modules + a footer, finally merged into a script file.

## Example

```
cat _header.tengo > my.tengo

cat block_google_secure_dns.tengo >> my.tengo
cat block_aaaa.tengo >> my.tengo

cat _footer.tengo >> my.tengo
```
