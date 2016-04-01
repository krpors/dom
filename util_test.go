package dom

import "testing"

// TestUtilEscape tests the convenience method to escape XML character data.
func TestUtilEscape(t *testing.T) {
	data := "<>&&amp;"
	str := escape(data)
	expected := "&lt;&gt;&amp;&amp;amp;"
	if str != expected {
		t.Logf("expected: %v", expected)
		t.Logf("acuta:    %v", str)
		t.Errorf("escaped sequence does not match")
	}
}
