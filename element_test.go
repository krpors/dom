package dom

import (
	"testing"
)

func TestElementGetters(t *testing.T) {
	root := newElement()
	root.SetTagName("pfx:rewt")

	if root.GetNodeName() != "pfx:rewt" {
		t.Errorf("node name is expected to be 'pfx:rewt', but was '%v'", root.GetNodeName())
	}
	if root.GetNodeType() != ElementNode {
		t.Errorf("elements are supposed to be of type '%v'", ElementNode)
	}
	if root.GetNodeValue() != "" {
		t.Errorf("node value of elements are not applicable and should therefore be empty")
	}
	if len(root.GetChildNodes()) != 0 {
		t.Errorf("initialized node list should be zero length")
	}
	if root.GetNodeName() != root.GetTagName() {
		t.Errorf("GetNodeName() should equal GetTagName(): %v != %v", root.GetNodeName(), root.GetTagName())
	}
	if root.GetLocalName() != "rewt" {
		t.Errorf("local name should be 'rewt', but was '%v'", root.GetLocalName())
	}
	if root.GetNamespacePrefix() != "pfx" {
		t.Errorf("namespace prefix should be 'pfx', but was '%v'", root.GetNamespacePrefix())
	}
	if root.GetAttribute("anything") != "" {
		t.Errorf("no attributes set, expected empty string, but got '%s'", root.GetAttribute("anything"))
	}
	// set that attribute, but find a non existant one. Must return empty as well.
	root.SetAttribute("anything", "goes")
	if root.GetAttribute("nonexistent") != "" {
		t.Error("expected empty string due to unfound attribute")
	}

	// add some children
	for i := 0; i < 10; i++ {
		e := newElement()
		e.SetTagName("element" + string(i))
		root.AppendChild(e)
	}
	if len(root.GetChildNodes()) != 10 {
		t.Errorf("node list length should be 10, but was %v", len(root.GetChildNodes()))
	}

	// element without prefix:
	root = newElement()
	root.SetTagName("rewt")
	if root.GetLocalName() != "rewt" {
		t.Errorf("local name should be 'rewt', was '%v'", root.GetLocalName())
	}
}

func TestElementOwnerDocument(t *testing.T) {
	doc := NewDocument()
	elem, err := doc.CreateElement("root")
	if err != nil {
		t.Errorf("no error expected!!")
	}
	if elem.GetOwnerDocument() != doc {
		t.Errorf("incorrect owner document")
	}
}

func TestElementAppendChild(t *testing.T) {
	root := newElement()
	root.SetTagName("root")

	child := newElement()
	child.SetTagName("child")

	if len(root.GetChildNodes()) != 0 {
		t.Errorf("length of node list should be 0, but was %v", len(root.GetChildNodes()))
	}

	if root.GetFirstChild() != nil {
		t.Error("first child should be nil")
	}

	err := root.AppendChild(child)
	if err != nil {
		t.Errorf("did not expect error at this point: '%v'", err)
	}

	if root.GetFirstChild() != child {
		t.Error("invalid first child")
	}

	if child.GetParentNode() != root {
		t.Logf("child.parent = %v, root = %v", child.GetParentNode(), root)
		t.Errorf("parent node of 'child' should be 'root'")
	}

	if len(root.GetChildNodes()) != 1 {
		t.Errorf("length of node list should be 1, but was %v", len(root.GetChildNodes()))
	}

	if root.GetChildNodes()[0] != child {
		t.Errorf("first element is expected to be 'child', but was %v", root.GetChildNodes()[0])
	}

	err = root.AppendChild(root)
	if err == nil {
		t.Errorf("expected a hierarchy error here")
	}

	attr := newAttr()
	err = root.AppendChild(attr)
	if err == nil {
		t.Error("expected error, got none")
	}
}

func TestElementHasChildNodes(t *testing.T) {
	root := newElement()
	if root.HasChildNodes() {
		t.Errorf("expected no child nodes")
	}

	// Adding root to itself should result in nothing
	root.AppendChild(root)
	if root.HasChildNodes() {
		t.Errorf("expected no child nodes")
	}

	child := newElement()
	root.AppendChild(child)
	if !root.HasChildNodes() {
		t.Errorf("expected child nodes")
	}
}

func TestElementAttributes(t *testing.T) {
	root := newElement()
	root.SetAttribute("cruft", "value")

	if root.GetAttribute("cruft") != "value" {
		t.Errorf("expected 'value', got '%v'", root.GetAttribute("cruft"))
	}

	attr := newAttr()
	attr.setName("pfx:anything")
	attr.setNamespaceURI("urn:any:namespace")
	attr.SetValue("harpy")

	root.SetAttributeNode(attr)

	if root.GetAttribute("pfx:anything") != "harpy" {
		t.Errorf("expected 'harpy' but was '%v'", root.GetAttribute("pfx:anything"))
	}
}

func TestElementGetElementsByTagName(t *testing.T) {
	root := newElement()
	root.SetTagName("root")

	child1 := newElement()
	child1.SetTagName("child")
	child1.SetAttribute("name", "a")

	child2 := newElement()
	child2.SetTagName("child")
	child2.SetAttribute("name", "b")

	child3 := newElement()
	child3.SetTagName("child")
	child3.SetAttribute("name", "ac")

	root.AppendChild(child1)
	child1.AppendChild(child3)
	root.AppendChild(child2)

	n1 := root.GetElementsByTagName("child")
	if len(n1) != 3 {
		t.Errorf("expected 3, got '%d'", len(n1))
		t.FailNow()
	}

	if n1[0].GetAttribute("name") != "a" {
		t.Errorf("expected 'a', got '%v'", n1[0].GetAttribute("name"))
	}
	if n1[1].GetAttribute("name") != "ac" {
		t.Errorf("expected 'ac', got '%v'", n1[1].GetAttribute("name"))
	}
	if n1[2].GetAttribute("name") != "b" {
		t.Errorf("expected 'b', got '%v'", n1[2].GetAttribute("name"))
	}

	n2 := child1.GetElementsByTagName("child")
	if len(n2) != 1 {
		t.Errorf("expected 1, got '%d'", len(n2))
	}

	// also, check equality of the pointers
	if n1[0] != child1 {
		t.Error("incorrect child node")
	}
	if n1[1] != child3 {
		t.Error("incorrect child node")
	}
	if n1[2] != child2 {
		t.Error("incorrect child node")
	}
}

func TestElementGetElementsByTagNameNS(t *testing.T) {
	// TODO: TestElementGetElementsByTagNameNS
}

// Tests the Lookup* methods on the Element type
func TestElementLookupNamespaceURI(t *testing.T) {
	doc := NewDocument()
	root, _ := doc.CreateElementNS("http://example.org/root", "pfx:root")
	root.SetAttribute("xmlns:abc", "http://example.org/cruft")

	parent, _ := doc.CreateElementNS("http://example.org/parent", "ns1:parent")

	child, _ := doc.CreateElement("child")
	child.SetAttribute("xmlns", "http://example.org/child")

	grandchild1, _ := doc.CreateElement("ns1:grandchild")

	grandchild2, _ := doc.CreateElement("othergrandchild")
	grandchild2.SetAttribute("xmlnsanythingafterthis", "http://example.org/grandchild2")

	doc.AppendChild(root)
	root.AppendChild(parent)
	parent.AppendChild(child)
	child.AppendChild(grandchild1)
	child.AppendChild(grandchild2)

	// Test lookup namespace stuff by URI:
	var tests = []struct {
		expected string
		actual   string
	}{
		{"http://example.org/root", root.LookupNamespaceURI("pfx")},
		{"http://example.org/child", child.LookupNamespaceURI("")},
		{"http://example.org/parent", grandchild1.LookupNamespaceURI("ns1")},
		{"http://example.org/root", grandchild1.LookupNamespaceURI("pfx")},
		{"http://example.org/grandchild2", grandchild2.LookupNamespaceURI("")},
		{"", grandchild2.LookupNamespaceURI("none-this-prefix-not-registered")},
	}

	for _, test := range tests {
		if test.expected != test.actual {
			t.Errorf("expected '%s', got '%s'", test.expected, test.actual)
		}
	}

	// Test lookup prefix stuff:
	tests = []struct {
		expected string
		actual   string
	}{
		{"pfx", root.LookupPrefix("http://example.org/root")},
		{"pfx", parent.LookupPrefix("http://example.org/root")},
		{"pfx", grandchild1.LookupPrefix("http://example.org/root")},
		{"ns1", parent.LookupPrefix("http://example.org/parent")},
		{"abc", grandchild1.LookupPrefix("http://example.org/cruft")},
		{"", grandchild1.LookupPrefix("urn:nonexistant:namespace")},
		{"", grandchild1.LookupPrefix("")},
	}

	for _, test := range tests {
		if test.expected != test.actual {
			t.Errorf("expected '%s', got '%s'", test.expected, test.actual)
		}
	}

}

func TestElementGetPreviousNextSibling(t *testing.T) {
	doc := NewDocument()
	root, _ := doc.CreateElement("root")
	first, _ := doc.CreateElement("first")
	second, _ := doc.CreateElement("second")
	third, _ := doc.CreateElement("third")
	fourth := doc.CreateText("some arbitrary text at the fourth position")
	fifth, _ := doc.CreateElement("fifth")

	noparent, _ := doc.CreateElement("noparent")

	doc.AppendChild(root)
	root.AppendChild(first)
	root.AppendChild(second)
	root.AppendChild(third)
	root.AppendChild(fourth)
	root.AppendChild(fifth)

	var tests = []struct {
		expected Node // or nil.
		actual   Node // or nil.
	}{
		{nil, first.GetPreviousSibling()},
		{first, second.GetPreviousSibling()},
		{fourth, fifth.GetPreviousSibling()},
		{nil, fifth.GetNextSibling()},
		{third, second.GetNextSibling()},
		{nil, root.GetPreviousSibling()},
		{nil, doc.GetPreviousSibling()},
		{nil, doc.GetNextSibling()},
		{nil, noparent.GetPreviousSibling()},
		{nil, noparent.GetNextSibling()},
	}

	for _, test := range tests {
		if test.actual != test.expected {
			t.Errorf("expected '%v', got '%v'", test.expected, test.actual)
		}
	}
}
