package filetype

import "testing"

func TestIsJSON(t *testing.T) {
	if !IsJSON(`{"name":"hello"}`) {
		t.Errorf("Should be json")
	}

	if !IsJSON(`[{"name":"hello"}]`) {
		t.Errorf("Should be json")
	}

	if IsJSON(`[{"name":"hello"}`) {
		t.Errorf("Should not be json")
	}
}
