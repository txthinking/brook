## Brook GUI

In Brook GUI, scripts are abstracted into modules, and it will automatically combine `_header.tengo` and `_footer.tengo`, so you only need to write the module itself.

```
modules = append(modules, {
    // If you want to predefine multiple brook links, and then programmatically specify which one to connect to, then define `brooklinks` key a function
    brooklinks: func(m) {
        // Please refer to the example in `brooklinks.tengo`
    },
    // If you want to intercept and handle a DNS query, then define `dnsquery` key a function, `m` is the `in_dnsquery`
    dnsquery: func(m) {
        // Please refer to the example in `block_aaaa.tengo`
    },
    // If you want to intercept and handle an address, then define `address` key a function, `m` is the `in_address`
    address: func(m) {
        // Please refer to the example in `block_google_secure_dns.tengo`
    },
    // If you want to intercept and handle a http request, then define `httprequest` key a function, `request` is the `in_httprequest`
    httprequest: func(request) {
        // Please refer to the example in `ios_app_downgrade.tengo` or `redirect_google_cn.tengo`
    },
    // If you want to intercept and handle a http response, then define `httpresponse` key a function, `request` is the `in_httprequest`, `response` is the `in_httpresponse`
    httpresponse: func(request, response) {
        // Please refer to the example in `response_sample.tengo`
    }
})
```

## tun2brook

If you are using tun2brook, you can combine multiple modules into a complete script in the following way. For example:

```
cat _header.tengo > my.tengo

cat block_google_secure_dns.tengo >> my.tengo
cat block_aaaa.tengo >> my.tengo

cat _footer.tengo >> my.tengo
```
