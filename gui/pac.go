package main

import "net/http"

func GetPAC(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-ns-proxy-autoconfig")
	w.Write([]byte("function FindProxyForURL(url, host){ return \"SOCKS5 local.txthinking.com:1080; SOCKS local.txthinking.com:1080; DIRECT\"; }"))
}
