// +build darwin
// +build amd64 386

package sysproxy

import "testing"

func TestGetNetworkServices(t *testing.T) {
	nss, err := getNetworkServices()
	if err != nil {
		t.Fatal(err)
	}
	for _, v := range nss {
		t.Log("|" + v + "|")
	}
}

func TestSetDNSServer(t *testing.T) {
	err := SetDNSServer("8.8.4.4")
	if err != nil {
		t.Fatal(err)
	}
}
