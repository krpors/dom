package dom

import (
	"testing"
)

func TestCreateElement(t *testing.T) {
	doc := NewDocument()
	elem, _ := doc.CreateElement("root_node")
	if elem.GetTagName() != "root_node" {
		t.Errorf("expected 'root_node'")
	}

	if elem.OwnerDocument() != doc {
		t.Errorf("incorrect owner document: got %v, want %v", elem.OwnerDocument(), doc)
	}
}

func TestCreateElementNS(t *testing.T) {
	doc := NewDocument()
	elem, _ := doc.CreateElementNS("http://example.org/2016/ns", "elem")
	if elem.NamespaceURI() != "http://example.org/2016/ns" {
		t.Errorf("got wrong namespace: %s", elem.NamespaceURI())
	}
}

func TestCreateTextNode(t *testing.T) {
	doc := NewDocument()
	text := doc.CreateTextNode("abc123")
	if text.OwnerDocument() != doc {
		t.Errorf("incorrect owner document: got %v, want %v", text.OwnerDocument(), doc)
	}

	if text.GetData() != "abc123" {
		t.Errorf("got %v, want %v", text.GetData(), "abc123")
	}
}

func TestGetDocumentElement(t *testing.T) {
}

// Tests the setting of attributes and that whole crapload.
func TestAttributes(t *testing.T) {
	doc := NewDocument()
	root, _ := doc.CreateElement("root")
	root.SetAttribute("first", "first value")
	root.SetAttribute("second", "second value")
	root.SetAttribute("third", "third value")
	root.SetAttribute("fourth", "fourth value")

	nnm := root.GetAttributes()
	if nnm.Length() != 4 {
		t.Errorf("expected 4 attributes, got ", root.GetAttributes().Length())
	}

	if nnm.GetNamedItem("first").NodeValue() != "first value" {
		t.Errorf("expected 'first value'")
	}

	if nnm.GetNamedItem("second").NodeValue() != "second value" {
		t.Errorf("expected 'second value'")
	}

	if nnm.GetNamedItem("third").NodeValue() != "third value" {
		t.Errorf("expected 'second value'")
	}

	if nnm.GetNamedItem("fourth").NodeValue() != "fourth value" {
		t.Errorf("expected 'second value'")
	}
}
