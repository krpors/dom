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

//
func TestUtilMoveNamespacesToRoot(t *testing.T) {
	doc := NewDocument()
	root, _ := doc.CreateElement("employees")
	root.SetAttribute("hasNamespace", "nope")

	employee1, _ := doc.CreateElementNS("urn:employee", "employee")
	person1, _ := doc.CreateElementNS("urn:person", "p1:person")
	person1.SetAttribute("name", "Mimi")
	person2, _ := doc.CreateElementNS("urn:person", "p2:person")
	person3, _ := doc.CreateElementNS("urn:person", "person")

	employee2, _ := doc.CreateElementNS("urn:employee", "employee")
	employee2.SetAttribute("nonamespace", "valuevalue")

	employee2ExtraInfo, _ := doc.CreateAttributeNS("urn:extraInfo", "pfx1:extraInfo")
	employee2ExtraInfo.SetValue("deb")
	employee2Bleh, _ := doc.CreateAttributeNS("urn:attr_no_pfx", "attributeNoPrefix")
	employee2Bleh.SetValue("HI THAR")

	elemExtraInfo, _ := doc.CreateElementNS("urn:extraInfo", "extraInfo")

	doc.AppendChild(root)

	root.AppendChild(employee1)
	root.AppendChild(employee2)
	root.AppendChild(elemExtraInfo)

	employee1.AppendChild(person1)
	employee1.AppendChild(person2)
	employee1.AppendChild(person3)

	employee2.SetAttributeNode(employee2ExtraInfo)
	employee2.SetAttributeNode(employee2Bleh)

	var b bytes.Buffer

	PrintTree(doc, &b)
	t.Logf("\nBefore moving namespaces to root:\n\n%s", b.String())

	// Move namespaces to the document element.
	MoveNamespacesToRoot(doc)

	b.Reset()
	PrintTree(doc, &b)
	t.Logf("\nAfter moving namespaces to root:\n\n%s", b.String())

	// Assert things.
	tests := []struct {
		expected interface{}
		actual   interface{}
	}{
		{"employees", root.GetTagName()},
		{6, root.GetAttributes().Length()},
		{"urn:employee", root.GetAttribute("xmlns:ns0")},
		{"urn:person", root.GetAttribute("xmlns:p1")},
		{"urn:person", root.GetAttribute("xmlns:p2")},
		{"urn:extraInfo", root.GetAttribute("xmlns:pfx1")},
		{"Mimi", person1.GetAttribute("name")},
		{"pfx1:extraInfo", root.GetChildNodes()[2].GetNodeName()},
	}

	for i, bla := range tests {
		if bla.expected != bla.actual {
			t.Errorf("test index %d: expected '%v', got '%v'", i, bla.expected, bla.actual)
		}
	}
}
