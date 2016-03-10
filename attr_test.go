package dom

import (
	"testing"
)

// Tests basic getters etc.
func TestAttrGetters(t *testing.T) {
	doc := NewDocument()
	elem, _ := doc.CreateElement("tag")
	a := newAttr()
	a.setOwnerElement(elem)
	a.setOwnerDocument(doc)
	a.setParentNode(newElement())
	a.setName("cruft")
	a.SetValue("valval")
	a.setNamespaceURI("http://example.org/lol")
	if a.GetName() != "cruft" {
		t.Errorf("incorrect node name")
	}
	if a.GetNodeName() != "cruft" {
		t.Errorf("incorrect node name")
	}
	if a.GetParentNode() != nil {
		t.Error("attr cannot have a parent (must be nil)")
	}
	if err := a.AppendChild(newElement()); err == nil {
		t.Error("expected an error at this point")
	}
	if len(a.GetChildNodes()) != 0 {
		t.Error("len of child nodes must be zero at all times")
	}
	if a.GetFirstChild() != nil {
		t.Error("first child must always be nil")
	}
	if a.GetAttributes() != nil {
		t.Error("attributes must always be nil")
	}
	if a.GetOwnerDocument() != doc {
		t.Error("incorrect owner document")
	}
	if a.HasChildNodes() != false {
		t.Error("must always return false, but was true")
	}
	if a.GetOwnerElement() != elem {
		t.Error("incorrect owner element")
	}
	if a.GetNodeType() != AttributeNode {
		t.Errorf("incorrect node type for attribute")
	}
	if a.GetNodeValue() != "valval" {
		t.Errorf("incorrect node value: '%v'", a.GetNodeValue())
	}
	if a.GetValue() != "valval" {
		t.Errorf("incorrect node value: '%v'", a.GetValue())
	}
}

func TestAttrLookupNamespaceURI(t *testing.T) {
	root := newElement()
	root.SetTagName("root")
	root.SetAttribute("xmlns:pfx", "http://example.org/pfx")
	root.SetAttribute("xmlns:xfb", "urn:xfbcft")

	child := newElement()
	child.SetTagName("child")
	child.SetAttribute("pfx:name", "Mimi")

	attr := child.GetAttributes().GetNamedItem("pfx:name")

	t.Logf("%s", attr.GetNodeName())
}
