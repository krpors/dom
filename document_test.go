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

func TestDocumentCreateElementInvalid(t *testing.T) {
	doc := NewDocument()
	elem, err := doc.CreateElement("in valid tag")
	if err == nil && elem != nil {
		t.Errorf("expected error due to invalid tag name")
	}
}
