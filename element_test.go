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
	if !root.HasAttributes() {
		t.Error("element should have attributes")
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

	// Try setting an attribute with an invalid name.
	err = root.SetAttribute("<invalid", "name")
	if err == nil {
		t.Error("expected error due to invalid attribute name")
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

	_, found = grandchild2.LookupNamespaceURI("none-this-prefix-is-not-registered")
	if found {
		t.Error("namespace prefix should not be found")
	}

	prefixLookupTest := func(src Node, namespaceToLookup string, expectedFound bool, expectedPrefix string) {
		actualPfx, actualFound := src.LookupPrefix(namespaceToLookup)
		if expectedFound != actualFound {
			t.Errorf("Expected to find prefix '%s' for namespace '%s', but was not found", expectedPrefix, namespaceToLookup)
		}

		if actualFound && actualPfx != expectedPrefix {
			t.Errorf("Prefix for namespace '%s' was resolved to '%s', but expected it to be '%s'", namespaceToLookup, actualPfx, expectedPrefix)
		}
	}

	// should be found
	prefixLookupTest(root, "http://example.org/root", true, "pfx")
	prefixLookupTest(parent, "http://example.org/root", true, "pfx")
	prefixLookupTest(grandchild1, "http://example.org/root", true, "pfx")
	prefixLookupTest(parent, "http://example.org/parent", true, "ns1")
	prefixLookupTest(grandchild1, "http://example.org/cruft", true, "abc")

	// should not be found
	prefixLookupTest(grandchild1, "urn:nonexistent:namespace", false, "")
	prefixLookupTest(grandchild1, "", false, "")
}

func TestElementIsDefaultNamespace1(t *testing.T) {
	doc := NewDocument()
	root, _ := doc.CreateElementNS("urn:root:ns", "root")
	if !root.IsDefaultNamespace("urn:root:ns") {
		t.Errorf("Expected the namespace to be the default namespace")
	}
}

func TestElementIsDefaultNamespace2(t *testing.T) {
	doc := NewDocument()
	root, _ := doc.CreateElementNS("urn:root:ns", "pfx:root")
	if root.IsDefaultNamespace("urn:root:ns") {
		t.Errorf("Expected the namespace to **NOT** be the default namespace")
	}
}

func TestElementIsDefaultNamespaceAttrTest(t *testing.T) {
	doc := NewDocument()
	root, _ := doc.CreateElement("pfx:root")
	root.SetAttribute("xmlns", "urn:root:ns")

	if !root.IsDefaultNamespace("urn:root:ns") {
		t.Errorf("Expected to find the default namespace")
	}
}

func TestElementIsDefaultNamespaceAncestors(t *testing.T) {
	doc := NewDocument()
	root, _ := doc.CreateElement("pfx:root")
	root.SetAttribute("xmlns", "urn:root:ns")

	child, _ := doc.CreateElement("pfx:somechild")
	root.AppendChild(child)

	if !child.IsDefaultNamespace("urn:root:ns") {
		// The ancestors should be visited
		t.Error("Expected to find the default namespace because it is declared in the root element")
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

// Tests InsertBefore functionality.
func TestElementInsertBefore(t *testing.T) {
	// <root>
	//   <child>
	//     <grandchild/>
	//   </child>
	// </root>
	doc := NewDocument()
	root, _ := doc.CreateElement("root")
	child, _ := doc.CreateElement("child")
	grandChild, _ := doc.CreateElement("grandchild")
	replacement, _ := doc.CreateElement("replacement")
	attr, _ := doc.CreateAttribute("attribute")

	doc.AppendChild(root)
	root.AppendChild(child)
	child.AppendChild(grandChild)

	if _, err := root.InsertBefore(nil, child); err == nil {
		t.Error("expected error (newChild is nil)")
	}
	if _, err := root.InsertBefore(attr, child); err == nil {
		t.Error("expected error (invalid nodetype replacement)")
	}

	// "Normal" replacement. refChild is found, same owner document etc.
	// The resulting document should now look like:
	// <root>
	//   <replacement/>
	//   <child>
	//     <grandchild/>
	//   </child>
	// </root>
	repl, err := root.InsertBefore(replacement, child)
	if err != nil {
		t.Error("unexpected error")
		t.FailNow()
	}
	if repl != replacement {
		t.Error("repl != replacement")
	}
	if len(root.GetChildNodes()) != 2 {
		t.Errorf("expected 2 child nodes in <root>, got %v", len(root.GetChildNodes()))
	}

	// Grandchild is already in the tree, must be removed.
	// The resulting document should now look like:
	// <root>
	//   <replacement/>
	//   <grandchild/>
	//   <child/>
	// </root>
	repl, err = root.InsertBefore(grandChild, child)
	if err != nil {
		t.Error("unexpected error")
		t.FailNow()
	}
	if repl != grandChild {
		t.Error("repl != grandchild")
	}

	// Double check the order of elements.
	if len(root.GetChildNodes()) != 3 {
		t.Errorf("expected 3 child nodes, got %v", len(root.GetChildNodes()))
	}

	cnode := root.GetChildNodes()[0]
	if cnode != replacement {
		t.Errorf("node 0 should be <replacement>, but was '%v'", cnode)
	}
	cnode = root.GetChildNodes()[1]
	if cnode != grandChild {
		t.Errorf("node 1 should be <grandchild>, but was '%v'", cnode)
	}
	cnode = root.GetChildNodes()[2]
	if cnode != child {
		t.Errorf("node 2 should be <child>, but was '%v'", cnode)
	}

	// Check if the <child> node doesn't have any children left (grandchild
	// is move to the root).
	if len(child.GetChildNodes()) != 0 {
		t.Error("<child> should have zero child nodes left")
	}
}

// Tests the other cases for InsertBefore.
func TestElementInsertBefore2(t *testing.T) {
	dom1 := NewDocument()
	root1, _ := dom1.CreateElement("root1")
	child1, _ := dom1.CreateElement("child1")
	grandchild1, _ := dom1.CreateElement("grandchild1")

	dom1.AppendChild(root1)
	root1.AppendChild(child1)
	child1.AppendChild(grandchild1)

	dom2 := NewDocument()
	root2, _ := dom2.CreateElement("root2")
	dom2.AppendChild(root2)

	// Insert element created from different document
	if _, err := root1.InsertBefore(root2, child1); err == nil {
		t.Error("expected an error")
	}

	// grandchild1 is no child of root1, so this should return ErrNotFound.
	if _, err := root1.InsertBefore(grandchild1, root1); err == nil {
		t.Error("expected an error")
	}

	// refChild is nil, should result in appending the child.
	n, err := root1.InsertBefore(grandchild1, nil)
	if err != nil {
		t.Error("unexpected error")
		t.FailNow()
	}
	if n != grandchild1 {
		t.Error("n != grandchild1")
		t.FailNow()
	}
	if len(root1.GetChildNodes()) != 2 {
		t.Error("!!!")
	}
	// First child of <root> should be <child1>, second <grandchild1>.
	if root1.GetChildNodes()[0] != child1 {
		t.Error("node 0 should be <child1>")
	}
	if root1.GetChildNodes()[1] != grandchild1 {
		t.Error("node 1 should be <grandchild>")
	}
}

// Tests the cloning process.
func TestElementClone(t *testing.T) {
	doc := NewDocument()
	root, _ := doc.CreateElement("prefix:Root")
	child, _ := doc.CreateElementNS("urn:uri", "ns:Child")
	child.SetAttribute("xmlns:attrpfx", "urn:attrpfx")
	child.SetAttribute("a", "b")
	attr, _ := doc.CreateAttributeNS("urn:attrpfx", "attrpfx:namespacedAttribute")
	attr.SetValue("example")

	doc.AppendChild(root)
	root.AppendChild(child)
	child.SetAttributeNode(attr)

	// Clone the child node and return it (no deep clone)
	clone := child.CloneNode(false)

	if clone == child {
		t.Error("unexpected equality (n == child)")
		t.FailNow()
	}

	tests := []struct {
		expected interface{}
		actual   interface{}
	}{
		{"urn:uri", clone.GetNamespaceURI()},
		{"ns:Child", clone.GetNodeName()},
		{3, clone.GetAttributes().Length()},
		{"b", clone.GetAttributes().GetNamedItem("a").GetNodeValue()},
		{"example", clone.GetAttributes().GetNamedItem("attrpfx:namespacedAttribute").GetNodeValue()},
		{"urn:attrpfx", clone.GetAttributes().GetNamedItem("xmlns:attrpfx").GetNodeValue()},
	}

	for _, test := range tests {
		if test.expected != test.actual {
			t.Errorf("expected '%v', got '%v'", test.expected, test.actual)
		}
	}

	// Attempt a deep clone, and append it to the child
	clone = root.CloneNode(true)
	child.AppendChild(clone)

	elems := doc.GetElementsByTagName("prefix:Root")
	if len(elems) != 2 {
		t.Errorf("expected length of 2, got %d", len(elems))
	}
	elems = doc.GetElementsByTagName("ns:Child")
	if len(elems) != 2 {
		t.Errorf("expected length of 2, got %d", len(elems))
	}
}

func TestElementImportNode(t *testing.T) {
	doc := NewDocument()
	root, _ := doc.CreateElement("root")
	root.SetAttribute("first", "root-1")
	child, _ := doc.CreateElement("child")
	child.SetAttribute("first", "child-1")
	child.SetAttribute("second", "child-2")
	grandchild, _ := doc.CreateElementNS("urn:grandchild:ns", "grandchild")
	grandchild.SetAttribute("first", "grandchild-1")
	grandchild.SetTextContent("hello")

	doc.AppendChild(root)
	root.AppendChild(child)
	child.AppendChild(grandchild)

	// Create a new dcument, import root node and all it's descendants.
	doc2 := NewDocument()
	imported := doc2.ImportNode(root, true)
	doc2.AppendChild(imported)

	if imported.GetOwnerDocument() != doc2 {
		t.Error("incorrect owner document")
	}

	// This testing table is a trainwreck, but it gets the job done.
	tests := []struct {
		expected interface{}
		actual   interface{}
	}{
		{1, len(doc2.GetChildNodes())},
		{"root", doc2.GetChildNodes()[0].GetNodeName()},
		{"root-1", doc2.GetChildNodes()[0].GetAttributes().GetNamedItem("first").GetNodeValue()},
		{"child", doc2.GetChildNodes()[0].GetChildNodes()[0].GetNodeName()},
		{"child-1", doc2.GetChildNodes()[0].GetChildNodes()[0].GetAttributes().GetNamedItem("first").GetNodeValue()},
		{"child-2", doc2.GetChildNodes()[0].GetChildNodes()[0].GetAttributes().GetNamedItem("second").GetNodeValue()},
		{"grandchild", doc2.GetChildNodes()[0].GetChildNodes()[0].GetChildNodes()[0].GetNodeName()},
		{"urn:grandchild:ns", doc2.GetChildNodes()[0].GetChildNodes()[0].GetChildNodes()[0].GetNamespaceURI()},
		{"hello", doc2.GetChildNodes()[0].GetChildNodes()[0].GetChildNodes()[0].GetChildNodes()[0].GetNodeValue()},
	}

	for i, test := range tests {
		if test.actual != test.expected {
			t.Errorf("test %d: expected %v, got %v", i, test.expected, test.actual)
			t.FailNow()
		}
	}
}
