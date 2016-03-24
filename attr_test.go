package dom

import (
	"testing"
)

// Tests basic getters etc.
func TestAttrGetters(t *testing.T) {
	doc := NewDocument()
	elem, _ := doc.CreateElement("tag")
	doc.AppendChild(elem)
	a, _ := doc.CreateAttributeNS("http://example.org/lol", "pfx:cruft")
	elem.SetAttributeNode(a)
	a.SetValue("valval")
	if a.GetName() != "pfx:cruft" {
		t.Error("incorrect node name")
	}
	if a.GetNodeName() != "pfx:cruft" {
		t.Error("incorrect node name")
	}
	if a.GetLocalName() != "cruft" {
		t.Error("incorrect node name")
	}
	if a.GetNamespacePrefix() != "pfx" {
		t.Error("incorrect prefix")
	}
	if a.GetNamespaceURI() != "http://example.org/lol" {
		t.Error("incorrect namespace URI")
	}
	if a.GetParentNode() != nil {
		t.Error("attr cannot have a parent (must be nil)")
	}
	bogusElem, _ := doc.CreateElement("bogus")
	if err := a.AppendChild(bogusElem); err == nil {
		t.Error("expected an error at this point")
	}
	a.setParentNode(bogusElem)
	if a.GetParentNode() != nil {
		t.Error("parent node should be nil at all times")
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
	if a.GetPreviousSibling() != nil {
		t.Error("expected nil previous sibling")
	}
	if a.GetNextSibling() != nil {
		t.Error("expected nil next sibling")
	}
	if a.GetLastChild() != nil {
		t.Error("expecting nil last child")
	}
}

func TestAttrLookupNamespaceURI(t *testing.T) {
	doc := NewDocument()

	root, _ := doc.CreateElement("root")
	root.SetAttribute("xmlns:pfx", "http://example.org/pfx")
	root.SetAttribute("xmlns:xfb", "urn:xfbcft")

	child, _ := doc.CreateElement("child")
	child.SetAttribute("pfx:name", "Mimi")

	attr := child.GetAttributes().GetNamedItem("pfx:name").(Attr)

	root.AppendChild(child)

	ns, found := attr.LookupNamespaceURI("pfx")
	exp := "http://example.org/pfx"
	if ns != exp || !found {
		t.Errorf("expected '%v', got '%v'", exp, ns)
	}

	// Attribute node owned by nothing:
	attr, _ = doc.CreateAttribute("no-owner")
	if _, found := attr.LookupNamespaceURI("pfxWhatever"); found {
		t.Error("expecting false")
	}
}

func TestAttrLookupPrefix(t *testing.T) {
	doc := NewDocument()
	root, _ := doc.CreateElementNS("urn:ns:attr1", "ns1:root")
	sub1, _ := doc.CreateElement("ns1:sub1")
	sub2, _ := doc.CreateElement("ns1:sub2")
	sub3, _ := doc.CreateElement("ns1:sub3")
	sub4, _ := doc.CreateElement("ns1:sub4")

	attr1, _ := doc.CreateAttribute("ns1:name")
	attr1.SetValue("melissandre")

	doc.AppendChild(root)
	root.AppendChild(sub1)
	root.AppendChild(sub2)
	root.AppendChild(sub3)
	root.AppendChild(sub4)
	sub4.SetAttributeNode(attr1)

	pfx := attr1.LookupPrefix("urn:ns:attr1")
	if pfx != "ns1" {
		t.Errorf("expected 'ns1', got '%v'", pfx)
	}

	// Attribute node owned by nothing:
	attr1, _ = doc.CreateAttribute("no-owner")
	if attr1.LookupPrefix("n") != "" {
		t.Error("expecting empty string")
	}
}

func TestAttrReplaceInsertRemoveChild(t *testing.T) {
	doc := NewDocument()
	attr, _ := doc.CreateAttribute("attr")
	if _, err := attr.ReplaceChild(nil, nil); err == nil {
		t.Error("expected error")
	}
	if _, err := attr.InsertBefore(nil, nil); err == nil {
		t.Error("expected error")
	}
	if _, err := attr.RemoveChild(nil); err == nil {
		t.Error("expected error")
	}
}
