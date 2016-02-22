package dom

import (
	"bytes"
	"testing"
)

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

func TestToXML(t *testing.T) {
	doc := NewDocument()

	root, _ := doc.CreateElement("root")
	first, _ := doc.CreateElement("first")
	text := doc.CreateTextNode("< & > are entities!")
	second, _ := doc.CreateElement("second")
	nonchildren, _ := doc.CreateElement("nochildren")

	doc.AppendChild(root)

	root.AppendChild(first)
	first.AppendChild(text)

	root.AppendChild(second)
	second.AppendChild(nonchildren)

	var buf bytes.Buffer
	ToXML(doc, false, &buf)

	expected := `<?xml version="1.0" encoding="UTF-8"?><root><first>&lt; &amp; &gt; are entities!</first><second><nochildren/></second></root>`

	if buf.String() != expected {
		t.Logf("actual:   %v", buf.String())
		t.Logf("expected: %v", expected)
		t.Errorf("unexpected output")
	}
}
