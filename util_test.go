package dom

import (
	"os"
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

func TestUtilMoveNamespacesToRoot(t *testing.T) {
	doc := NewDocument()
	root, _ := doc.CreateElement("root-no-namespace")
	root.SetAttribute("xmlns:cruft", "hi")
	root.SetAttribute("xmlns", "default")
	root.SetAttribute("lunis", "vortalds")
	child, _ := doc.CreateElementNS("urn:child", "prefix:child")
	child.SetAttribute("xmlns", "childnamespace")
	subchild, _ := doc.CreateElement("none:hi")

	doc.AppendChild(root)
	root.AppendChild(child)
	child.AppendChild(subchild)

	MoveNamespacesToRoot(doc)
	// doc.NormalizeDocument()
	PrintTree(doc, os.Stdout)
}
