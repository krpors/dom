package dom

import (
	"fmt"
	"testing"
)

func TestElementGetters(t *testing.T) {
	doc := NewDocument()
	root, _ := doc.CreateElement("pfx:rewt")

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
		e, _ := doc.CreateElement(fmt.Sprintf("element%d", i))
		root.AppendChild(e)
	}
	if len(root.GetChildNodes()) != 10 {
		t.Errorf("node list length should be 10, but was %v", len(root.GetChildNodes()))
	}

	child := root.GetFirstChild().GetNodeName()
	if child != "element0" {
		t.Errorf("first child should be 'element0', but was '%v'", child)
	}
	child = root.GetLastChild().GetNodeName()
	if child != "element9" {
		t.Errorf("first child should be 'element9', but was '%v'", child)
	}

	// element without prefix:
	root, _ = doc.CreateElement("rewt")
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
	doc := NewDocument()
	root, _ := doc.CreateElement("root")
	child, _ := doc.CreateElement("child")

	if len(root.GetChildNodes()) != 0 {
		t.Errorf("length of node list should be 0, but was %v", len(root.GetChildNodes()))
	}

	if root.GetFirstChild() != nil {
		t.Error("first child should be nil")
	}

	if root.GetLastChild() != nil {
		t.Error("last child should be nil")
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

	attr, _ := doc.CreateAttribute("attr")
	err = root.AppendChild(attr)
	if err == nil {
		t.Error("expected error, got none")
	}
}

func TestElementHasChildNodes(t *testing.T) {
	doc := NewDocument()
	root, _ := doc.CreateElement("tag")
	if root.HasChildNodes() {
		t.Errorf("expected no child nodes")
	}

	// Adding root to itself should result in nothing
	root.AppendChild(root)
	if root.HasChildNodes() {
		t.Errorf("expected no child nodes")
	}

	child, _ := doc.CreateElement("child")
	root.AppendChild(child)
	if !root.HasChildNodes() {
		t.Errorf("expected child nodes")
	}
}

func TestElementAttributes(t *testing.T) {
	doc := NewDocument()
	root, _ := doc.CreateElement("root")
	// Declare one namespace for usage.
	root.SetAttribute("xmlns:declared", "http://example.org/declared")

	root.SetAttribute("cruft", "value")

	if root.GetAttribute("cruft") != "value" {
		t.Errorf("expected 'value', got '%v'", root.GetAttribute("cruft"))
	}

	attr, _ := doc.CreateAttributeNS("urn:any:namespace", "pfx:anything")
	attr.SetValue("harpy")
	root.SetAttributeNode(attr)

	if root.GetAttribute("pfx:anything") != "harpy" {
		t.Errorf("expected 'harpy' but was '%v'", root.GetAttribute("pfx:anything"))
	}

	// Setting an attribute with a prefix, but a namespace cannot be found
	// should generate an error.
	err := root.SetAttribute("undeclared:attr", "fail")
	if err == nil {
		t.Error("expected an error")
	}

	// Setting an attribute, prefix is declared in the root. Should be OK.
	err = root.SetAttribute("declared:name", "captain planet")
	if err != nil {
		t.Error("unexpected error")
	}
}

func TestElementSetAttributeNodeWrongDoc(t *testing.T) {
	docOne := NewDocument()
	e, _ := docOne.CreateElement("rootOne")

	docTwo := NewDocument()
	attr, _ := docTwo.CreateAttribute("attr")

	// Attempt to set an attribute created by docTwo, into docOne.
	if err := e.SetAttributeNode(attr); err == nil {
		t.Error("expected a error")
	}
}

func TestElementSetAttributeNodeOwnedByElement(t *testing.T) {
	doc := NewDocument()
	e1, _ := doc.CreateElement("root")
	e2, _ := doc.CreateElement("sub")
	a, _ := doc.CreateAttribute("attr")

	doc.AppendChild(e1)
	e1.AppendChild(e2)
	if err := e1.SetAttributeNode(a); err != nil {
		t.Error("unexpected error")
	}

	if err := e2.SetAttributeNode(a); err == nil {
		t.Error("expected error")
	}
}

func TestElementGetElementsByTagName(t *testing.T) {
	doc := NewDocument()
	root, _ := doc.CreateElement("root")

	child1, _ := doc.CreateElement("child")
	child1.SetAttribute("name", "a")

	child2, _ := doc.CreateElement("child")
	child2.SetAttribute("name", "b")

	child3, _ := doc.CreateElement("child")
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
	expected := "http://example.org/root"
	actual, found := root.LookupNamespaceURI("pfx")
	if !found {
		t.Error("namespace prefix should have been found")
	} else {
		if actual != expected {
			t.Errorf("expected '%v', got '%v'", expected, actual)
		}
	}

	expected = "http://example.org/child"
	actual, found = child.LookupNamespaceURI("")
	if !found {
		t.Error("empty namespace prefix should have been found")
	} else {
		if actual != expected {
			t.Errorf("expected '%v', got '%v'", expected, actual)
		}
	}

	expected = "http://example.org/parent"
	actual, found = grandchild1.LookupNamespaceURI("ns1")
	if !found {
		t.Error("empty namespace prefix should have been found")
	} else {
		if actual != expected {
			t.Errorf("expected '%v', got '%v'", expected, actual)
		}
	}

	expected = "http://example.org/root"
	actual, found = grandchild1.LookupNamespaceURI("pfx")
	if !found {
		t.Error("empty namespace prefix should have been found")
	} else {
		if actual != expected {
			t.Errorf("expected '%v', got '%v'", expected, actual)
		}
	}

	expected = "http://example.org/grandchild2"
	actual, found = grandchild2.LookupNamespaceURI("")
	if !found {
		t.Error("empty namespace prefix should have been found")
	} else {
		if actual != expected {
			t.Errorf("expected '%v', got '%v'", expected, actual)
		}
	}

	expected = ""
	actual, found = grandchild2.LookupNamespaceURI("none-this-prefix-is-not-registered")
	if found {
		t.Error("namespace prefix should not be found")
	}

	// Test lookup prefix stuff:
	var tests = []struct {
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

func TestElementRemoveChild(t *testing.T) {
	doc := NewDocument()
	root, _ := doc.CreateElement("root")
	first, _ := doc.CreateElement("first")
	second, _ := doc.CreateElement("second")
	third, _ := doc.CreateElement("third")
	childOfThird, _ := doc.CreateElement("childOfThird")

	doc.AppendChild(root)
	root.AppendChild(first)
	root.AppendChild(second)
	root.AppendChild(third)
	third.AppendChild(childOfThird)

	nodeCount := len(root.GetChildNodes())
	if nodeCount != 3 {
		t.Errorf("expected 3 child nodes, got %v", nodeCount)
	}

	removedNode, err := root.RemoveChild(second)
	if err != nil {
		t.Error("unexpected error")
		t.FailNow()
	}

	if removedNode != second {
		t.Error("removedNode should equal second")
		t.FailNow()
	}

	// After removing, we expect 2 children: first and third
	nodeCount = len(root.GetChildNodes())
	if nodeCount != 2 {
		t.Errorf("expected 2 child nodes, got %v", nodeCount)
		t.FailNow()
	}

	if root.GetChildNodes()[0] != first && root.GetChildNodes()[1] != third {
		t.Error("incorrect children found")
		t.FailNow()
	}

	// Attempt to remove a child from root which is not a child.
	removedNode, err = root.RemoveChild(childOfThird)
	if err == nil {
		t.Error("expected an error at this point, but got none")
	}

	// Removing nil:
	removedNode, err = root.RemoveChild(nil)
	if removedNode != nil && err != nil {
		t.Error("removedNode and err should be nil")
	}
}

func TestElementReplaceNode(t *testing.T) {
	doc := NewDocument()
	root, _ := doc.CreateElement("root")
	first, _ := doc.CreateElement("first")
	second, _ := doc.CreateElement("second")
	third, _ := doc.CreateElement("third")
	childOfThird, _ := doc.CreateElement("childOfThird")

	// Replacement new element, with no parent yet.
	replacement, _ := doc.CreateElement("replacement")

	doc.AppendChild(root)
	root.AppendChild(first)
	root.AppendChild(second)
	root.AppendChild(third)
	third.AppendChild(childOfThird)

	// Replace the second with replacement
	theNode, err := root.ReplaceChild(replacement, second)
	if err != nil {
		t.Error("unexpected error")
		t.FailNow()
	}

	if theNode != second {
		t.Error("theNode should equal second")
		t.FailNow()
	}

	if len(root.GetChildNodes()) != 3 {
		t.Errorf("expected 3 child nodes, got %v", len(root.GetChildNodes()))
	}

	if root.GetChildNodes()[1] != replacement {
		t.Errorf("node at index 1 should be 'replacement', but was %v", root.GetChildNodes()[1])
	}

	if replacement.GetParentNode() != root {
		t.Error("incorrect parent node")
	}
}

func TestElementTextContent(t *testing.T) {
	doc := NewDocument()
	root, _ := doc.CreateElement("root")

	a, _ := doc.CreateElement("a")
	aText := doc.CreateText("hello")
	aComment, _ := doc.CreateComment("this comment should be ignored")

	b, _ := doc.CreateElement("b")
	bText := doc.CreateText("world")

	c, _ := doc.CreateElement("c")
	cText := doc.CreateText("watching")

	ba, _ := doc.CreateElement("ba")
	baText1 := doc.CreateText("thanks")
	baText2 := doc.CreateText("for")

	// Create the tree
	doc.AppendChild(root)

	// No children, return empty string.
	if root.GetTextContent() != "" {
		t.Error("root has no children, should return empty string")
	}

	root.AppendChild(a)
	root.AppendChild(b)
	root.AppendChild(c)

	a.AppendChild(aText)    // hello
	a.AppendChild(aComment) // <!-- this comment should be ignored -->

	b.AppendChild(bText) // world
	b.AppendChild(ba)

	ba.AppendChild(baText1) // thanks
	ba.AppendChild(baText2) // for

	c.AppendChild(cText) // watching

	txt := root.GetTextContent()
	expected := "helloworldthanksforwatching"
	if txt != expected {
		t.Errorf("expected '%v', got '%v'", expected, txt)
	}

	// Now, remove them all by calling SetTextContent:
	root.SetTextContent("HAI!")
	root.SetTextContent("") // this should do absolutely nothing.

	if len(root.GetChildNodes()) != 1 {
		t.Error("expected 1 child node after SetTextContent")
		t.FailNow()
	}

	if txt, ok := root.GetChildNodes()[0].(Text); ok {
		if txt.GetNodeValue() != "HAI!" {
			t.Error("incorrect node content")
		}
	} else {
		t.Error("failed type assertion for Text node")
	}
}
