// +build darwin
// +build amd64 386

package sysproxy

import (
	"log"
	"testing"
)

func TestGetNetworkService(t *testing.T) {
	log.Println(GetNetworkInterface())
}
