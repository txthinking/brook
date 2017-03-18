package brook

import "testing"

func TestMakeSSKey(t *testing.T) {
	key := MakeSSKey("a")
	t.Log(len(key))
	t.Log(key)
}
