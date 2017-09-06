package cmdline

import "testing"

func TestArgs0(t *testing.T) {
	if Args0() == "" {
		t.Fatal("expected a name")
	}
}
