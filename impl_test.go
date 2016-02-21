package dom

import (
	"testing"
)

/*
TODO: this file has tests in a pretty much illogical structure.... groups things in a more
orderly fashion cus this sucks ass.
*/

// Tests that the parent nodes are correctly set.
func TestParentNode(t *testing.T) {
	doc := NewDocument()
	if doc.ParentNode() != nil {
		t.Errorf("document should not have a parent node")
	}

	root, _ := doc.CreateElement("root")
	if root.ParentNode() != doc {
		t.Errorf("root should have the document as parent node")
	}

	sub1, _ := doc.CreateElement("sub_element")
	t.Log(sub1)
	root.AppendChild(sub1)
	if sub1.ParentNode() != root {
		t.Errorf("sub_element should have root as parent node, but was %v", sub1.ParentNode())
	}
}

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

	if elem.OwnerDocument() != doc {
		t.Errorf("incorrect owner document: got %v, want %v", elem.OwnerDocument(), doc)
	}

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

func TestCreateAttrNode(t *testing.T) {
	doc := NewDocument()
	attr, err := doc.CreateAttribute("attribute")
	if err != nil {
		t.Error("did not expect an error at this point")
	}
	if attr.OwnerDocument() != doc {
		t.Errorf("incorrect owner document: got %v, want %v", attr.OwnerDocument(), doc)
	}
}

func TestGetDocumentElement(t *testing.T) {
}

// Tests the Attr interface implementation.
func TestAttrNameAndValue(t *testing.T) {
	attr := newAttr("attribute")
	attr.SetValue("somevalue")

	// NodeName() and GetName(), and NodeValue() and GetValue() are supposed
	// to be identical, as per the spec.

	if attr.NodeName() != "attribute" {
		t.Errorf("expected '%s', was '%s'", "attribute", attr.NodeName())
	}
	if attr.GetName() != "attribute" {
		t.Errorf("expected '%s', was '%s'", "attribute", attr.GetName())
	}

	if attr.NodeValue() != "somevalue" {
		t.Errorf("expected '%s', was '%s'", "somevalue", attr.NodeValue())
	}

	if attr.GetValue() != "somevalue" {
		t.Errorf("expected '%s', was '%s'", "somevalue", attr.GetValue())
	}
}

func TestNamedNodeMapSetNamedItem(t *testing.T) {
	nnm := newNamedNodeMap()
	attr := newAttr("someattr")
	err := nnm.SetNamedItem(attr)
	if err != nil {
		t.Errorf("did not expect an error here")
	}

	element := newElement()
	err = nnm.SetNamedItem(element)
	if err == nil {
		t.Errorf("expected hierarchy request error, got nil")
	}
	text := newText()
	err = nnm.SetNamedItem(text)
	if err == nil {
		t.Errorf("expected hierarchy request error, got nil")
	}
}

// Tests the setting of attributes and that whole crapload.
func TestElementSetAttributes(t *testing.T) {
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
