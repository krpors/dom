package dom

import (
	"testing"
)

func TestDocumentNodeName(t *testing.T) {
	doc := NewDocument()
	if doc.NodeName() != "#document" {
		t.Errorf("node name of document must be '#document'")
	}
}

func TestDocumentNodeType(t *testing.T) {
	doc := NewDocument()
	if doc.NodeType() != DocumentNode {
		t.Errorf("node type of document must be %v", DocumentNode)
	}
}

func TestDocumentAppendChild(t *testing.T) {
	doc := NewDocument()

	err := doc.AppendChild(doc)
	if err == nil {
		t.Errorf("expected hierarchy error")
	}

	text := doc.CreateTextNode("HAI!")
	err = doc.AppendChild(text)
	if err == nil {
		t.Errorf("expected error at this point")
	}

	elem1, _ := doc.CreateElement("elem1")
	err = doc.AppendChild(elem1)
	if err != nil {
		t.Errorf("unexpected error")
	}

	elem2, _ := doc.CreateElement("elem2")
	err = doc.AppendChild(elem2)
	if err == nil {
		t.Errorf("expected error due to adding two root nodes")
	}
}

func TestDocumentHasChildNodes(t *testing.T) {
	doc := NewDocument()
	if doc.HasChildNodes() {
		t.Errorf("did not expect child nodes at this point")
	}

	elem, _ := doc.CreateElement("egregrious")
	err := doc.AppendChild(elem)
	if err != nil {
		t.Errorf("unexpected error")
	}

	if !doc.HasChildNodes() {
		t.Errorf("expected HasChildNodes to be true")
	}
}

// TestDocumentCreateelement tests the creation of elements using valid and
// invalid names. TODO: needs work. The isNameString() I borrowed fails with
// start characters which are supposed to be correct according to the spec.
// Or I am screwing things up?
func TestDocumentCreateElement(t *testing.T) {
	var tests = []struct {
		element        string
		expectedToBeOK bool
	}{
		{"valid", true},
		{"cruft_a", true},
		{"in val id", false},
		{"hi", true},
		{"  test", false},
		{":cruft", true},
		{"_ALAKAZAM", true},
		{":_something0123Darkside", true},
		{"øøøølmo", true},
		{"Grøups", true},
		{"\xc3\xb8stuff", true},
		{"...element", false},
		{"element...", true},
		{"Ållerskåléèöí", true},
		{"_More.Stuff.InXML.", true},
	}

	doc := NewDocument()
	for _, test := range tests {
		_, err := doc.CreateElement(test.element)
		if test.expectedToBeOK && err != nil {
			t.Errorf("XML name '%v' should be valid, but returned an error", test.element)
		} else if !test.expectedToBeOK && err == nil {
			t.Errorf("XML name '%v' should return an error, but was valid", test.element)
		}
	}
}

// ø
