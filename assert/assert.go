package assert

import "testing"

// Ok asserts if the boolean is true or false
func Ok(t *testing.T, b bool) {
	if b == false {
		t.Fatal("not ok")
	}
}
