package main

import "net/http"

func GetPAC(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-ns-proxy-autoconfig")
	s := `
function FindProxyForURL(url, host){
	// internal
    if(/\d+\.\d+\.\d+\.\d+/.test(host)){
        if (isInNet(dnsResolve(host), "10.0.0.0", "255.0.0.0") ||
                isInNet(dnsResolve(host), "172.16.0.0",  "255.240.0.0") ||
                isInNet(dnsResolve(host), "192.168.0.0", "255.255.0.0") ||
                isInNet(dnsResolve(host), "127.0.0.0", "255.255.255.0")){
            return "DIRECT";
        }
    }

	// plain
    if (isPlainHostName(host)){
        return "DIRECT";
    }

	if(dnsDomainIs(host, "txthinking.com")){
		return "DIRECT";
	}
	return "SOCKS5 local.txthinking.com:1080; SOCKS local.txthinking.com:1080; DIRECT";
}
`
	w.Write([]byte(s))
}
