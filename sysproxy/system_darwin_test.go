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
