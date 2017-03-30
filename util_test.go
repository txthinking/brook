package brook

import "testing"

func TestHasPort(t *testing.T) {
	var ok bool
	ok = HasPort("www.txthinking.com")
	if ok {
		t.Fatal("0")
	}
	ok = HasPort("www.txthinking.com:443")
	if !ok {
		t.Fatal("1")
	}
	ok = HasPort("[::1]")
	if ok {
		t.Fatal("2")
	}
	ok = HasPort("[::1]:80")
	if !ok {
		t.Fatal("3")
	}
}
