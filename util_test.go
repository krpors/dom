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

//
func TestUtilMoveNamespacesToRoot(t *testing.T) {
	doc := NewDocument()
	root, _ := doc.CreateElement("employees")
	root.SetAttribute("hasNamespace", "nope")

	employee1, _ := doc.CreateElementNS("urn:employee", "employee")
	person1, _ := doc.CreateElementNS("urn:person", "p1:person")
	person1.SetAttribute("name", "Mimi")
	person2, _ := doc.CreateElementNS("urn:person", "p2:person")
	person3, _ := doc.CreateAttributeNS("urn:person", "person")

	employee2, _ := doc.CreateElementNS("urn:employee", "employee")

	doc.AppendChild(root)

	root.AppendChild(employee1)
	root.AppendChild(employee2)

	employee1.AppendChild(person1)
	employee1.AppendChild(person2)
	employee1.AppendChild(person3)

	// Move namespaces to the document element.
	MoveNamespacesToRoot(doc)

	// Assert things.
	tests := []struct {
		expected interface{}
		actual   interface{}
	}{
		{"employees", root.GetTagName()},
		{4, root.GetAttributes().Length()},
		{"urn:employee", root.GetAttribute("xmlns:ns0")},
		{"urn:person", root.GetAttribute("xmlns:p1")},
		{"urn:person", root.GetAttribute("xmlns:p2")},
		{"Mimi", person1.GetAttribute("name")},
	}

	for i, bla := range tests {
		if bla.expected != bla.actual {
			t.Errorf("test index %d: expected '%v', got '%v'", i, bla.expected, bla.actual)
		}
	}
}
