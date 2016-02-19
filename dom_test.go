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
	if elem.NamespaceUri() != "http://example.org/2016/ns" {
		t.Errorf("got wrong namespace: %s", elem.NamespaceUri())
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
