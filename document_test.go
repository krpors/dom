package dom

import "testing"

// Test the plain getters of the Document. Also some no-op setters.
func TestDocumentGetters(t *testing.T) {
	doc := NewDocument()

	// Set some no-op setters, to increase coverage reports as well.
	doc.setParentNode(doc.CreateText("no-op"))

	if doc.GetNodeName() != "#document" {
		t.Errorf("node name of document must be '#document'")
	}
	if doc.GetLocalName() != "" {
		t.Error("local name must be an empty string")
	}
	if doc.GetNodeType() != DocumentNode {
		t.Errorf("node type of document must be %v", DocumentNode)
	}
	if doc.GetNodeValue() != "" {
		t.Error("node value must be an empty string")
	}
	if doc.GetAttributes() != nil {
		t.Error("attributes must be nil")
	}
	if doc.GetOwnerDocument() != nil {
		t.Error("owner document must be nil")
	}
	if doc.GetParentNode() != nil {
		t.Error("document can't have a parent node")
	}
	if doc.GetNamespaceURI() != "" {
		t.Error("document namespace uri should always be an empty string")
	}
	if doc.GetNamespacePrefix() != "" {
		t.Error("document cannot have a namespace prefix")
	}
}

// Tests appending doc to itself, and appending two document elements.
func TestDocumentAppendChild(t *testing.T) {
	doc := NewDocument()

	err := doc.AppendChild(nil)
	if err != nil {
		t.Error("appending a nil child should not generate an error")
	}

	err = doc.AppendChild(doc)
	if err == nil {
		t.Errorf("expected hierarchy error")
	}

	text := doc.CreateText("HAI!")
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

// Tests the appending of invalid children to a Document.
func TestDocumentAppendInvalidChildren(t *testing.T) {
	doc := NewDocument()
	text := doc.CreateText("should fail")
	if doc.AppendChild(text) == nil {
		t.Error("appending a text node to a document should fail")
	}
	attr, err := doc.CreateAttribute("hi")
	if err != nil {
		t.Error("??")
	}
	if doc.AppendChild(attr) == nil {
		t.Error("appending an attr to a document should fail")
	}
}

// Tests whether inserting processing instructions just works. They
// can appear anywhere in the document, before or after the root node.
func TestDocumentAppendChildProcInst(t *testing.T) {
	doc := NewDocument()

	elemRoot, _ := doc.CreateElement("root")
	procInst, _ := doc.CreateProcessingInstruction("lom", "lobon")
	procInst2, _ := doc.CreateProcessingInstruction("dowan", "duvessa")

	// Append it before anything else:
	err := doc.AppendChild(procInst)
	if err != nil {
		t.Error("unexpected error while adding process instruction")
	}

	if len(doc.GetChildNodes()) != 1 {
		t.Error("expected 1 child node at this point")
	}

	// Append the root node.
	err = doc.AppendChild(elemRoot)
	if err != nil {
		t.Error("unexpected error while adding element")
	}

	if len(doc.GetChildNodes()) != 2 {
		t.Error("expected 2 child nodes at this point")
	}

	// Append another processing instruction.
	err = doc.AppendChild(procInst2)
	if err != nil {
		t.Error("unexpected error while adding process instruction")
	}

	if len(doc.GetChildNodes()) != 3 {
		t.Error("expected 3 child nodes at this point")
	}
}

// Tests the functionality of inserting/appending elements with a different document owner.
func TestDocumentInvalidOwner(t *testing.T) {
	doc1 := NewDocument()
	doc2 := NewDocument()

	e, _ := doc2.CreateElement("doc2element")

	err := doc1.AppendChild(e)
	if err == nil {
		t.Error("expected an error, got none")
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
// invalid names, according to the XML spec.
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

// Tests the creation of elements with namespaces uris.
func TestDocumentCreateElementNS(t *testing.T) {
	doc := NewDocument()
	var tests = []struct {
		element        string
		namespace      string
		expectedToBeOK bool
	}{
		{"valid", "http://example.org/uri", true},
		{"valid", "uri:urn:bleh", true},
		{"cruft:valid", "anything can be put in this namespace", true},
		{"cruft:valid", "even w³ird, chøract€r$", true},
		{":zoit", "the XML spec doesn't care", true},
		{"¼ofanelement", "gopher://meh", false},
	}

	for _, test := range tests {
		_, err := doc.CreateElementNS(test.namespace, test.element)
		if test.expectedToBeOK && err != nil {
			t.Errorf("XML Name '%v' with namespace '%v' should be valid, but returned an error", test.element, test.namespace)
		} else if !test.expectedToBeOK && err == nil {
			t.Errorf("XML Name '%v' with namespace '%v' should return an error, but was valid", test.element, test.namespace)
		}
	}
}

func TestDocumentCreateComment(t *testing.T) {
	doc := NewDocument()
	c, err := doc.CreateComment("<anything goes in comments>")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if c.GetComment() != "<anything goes in comments>" {
		t.Error("incorrect comment content")
	}

	_, err = doc.CreateComment("except -- in comments")
	if err == nil {
		t.Error("expected an error during comment creation but got none")
	}
}

func TestDocumentCreateAttributeInvalid(t *testing.T) {
	doc := NewDocument()
	_, err := doc.CreateAttributeNS("urn:whatevs", "")
	if err == nil {
		t.Error("expected error at this point, but got none")
	}
}

func TestDocumentCreateAttributeNS(t *testing.T) {
	doc := NewDocument()
	root, err := doc.CreateElementNS("http://example.org/uri", "root")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	attr, err := doc.CreateAttributeNS("http://example.org/uri", "uri:name")
	attr.SetValue("zelda")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	doc.AppendChild(root)
	root.SetAttributeNode(attr)
}

func TestDocumentGetElementsBy(t *testing.T) {
	doc := NewDocument()
	root, _ := doc.CreateElement("root")
	child1, _ := doc.CreateElement("child")
	child2, _ := doc.CreateElement("child")
	child3, _ := doc.CreateElement("child")
	child4, _ := doc.CreateElementNS("http://example.org/ns1", "ns1:child")

	doc.AppendChild(root)
	root.AppendChild(child1)
	child1.AppendChild(child3)
	root.AppendChild(child2)
	root.AppendChild(child4)

	elems := doc.GetElementsByTagName("child")
	if len(elems) != 3 {
		t.Errorf("expected 3 elements, but got '%v'", len(elems))
	}

	elems = doc.GetElementsByTagNameNS("http://example.org/ns1", "child")
	if len(elems) != 1 {
		t.Errorf("expected 1 element, but got %d", len(elems))
	}
}

func TestDocumentInsertBefore(t *testing.T) {
	doc := NewDocument()
	pi, _ := doc.CreateProcessingInstruction("quux", "foo")
	doc.AppendChild(pi)

	// Result after AppendChild:
	// <document>
	//    <?quux foo?>

	root, _ := doc.CreateElement("rewt")
	n, err := doc.InsertBefore(root, pi)

	// Result of InsertBefore:
	// <document>
	//     <root/>
	//     <?quux foo?>

	if err != nil {
		t.Errorf("unexpected error: %v", err)
		t.FailNow()
	}

	if n != root {
		t.Error("returned node should be the same instance as the created node")
		t.FailNow()
	}

	if len(doc.GetChildNodes()) != 2 {
		t.Error("expected 2 child nodes")
		t.FailNow()
	}

	if doc.GetChildNodes()[0] != root {
		t.Error("incorrect first child node")
	}

	if doc.GetChildNodes()[1] != pi {
		t.Error("incorrect second child node")
	}

	// Check inserting a processing instruction already in the document
	doc.InsertBefore(pi, root)
	// Result after InsertBefore should be:
	// <document>
	//    <?quux foo?>
	//    <root>
	if len(doc.GetChildNodes()) != 2 {
		t.Error("expected 2 child nodes")
		t.FailNow()
	}
	if doc.GetChildNodes()[0] != pi {
		t.Error("incorrect first child node")
	}
	if doc.GetChildNodes()[1] != root {
		t.Error("incorrect second child node")
	}
	if pi.GetParentNode() != doc {
		t.Error("processing instruction has wrong parent node")
	}
	if root.GetParentNode() != doc {
		t.Error("root has wrong parent node")
	}

	// Nil new child should generate an error
	if _, err := doc.InsertBefore(nil, pi); err == nil {
		t.Error("expected an error")
	}

	// Adding another element should fail due to existing document element.
	e, _ := doc.CreateElement("another")
	if _, err = doc.InsertBefore(e, pi); err == nil {
		t.Error("expected an error")
	}

	attr, _ := doc.CreateAttribute("attr")
	text := doc.CreateText("text")
	if _, err = doc.InsertBefore(attr, pi); err == nil {
		t.Error("expected an error")
	}
	if _, err = doc.InsertBefore(text, pi); err == nil {
		t.Error("expected an error")
	}

	// Test inserting an element created from a different document.
	doc2 := NewDocument()
	pi2, _ := doc2.CreateProcessingInstruction("a", "b")
	if _, err = doc.InsertBefore(pi2, pi); err == nil {
		t.Error("expected an error")
	}

	// Test that inserting an element which is not a child returns a not found error.
	element, _ := doc.CreateElement("anything")
	if _, err = doc.InsertBefore(pi, element); err == nil {
		t.Error("expected an error")
	}
}

func TestDocumentRemoveChild(t *testing.T) {
	doc := NewDocument()
	e, _ := doc.CreateElement("root")
	if _, err := doc.RemoveChild(e); err == nil {
		t.Error("expected an error")
	}

	if a, b := doc.RemoveChild(nil); a != nil && b != nil {
		t.Error("returned Node and error should both be nil")
	}
}

func TestDocumentReplaceChild(t *testing.T) {
	doc := NewDocument()
	root, _ := doc.CreateElement("root")
	childOfRoot, _ := doc.CreateElement("child_of_root")
	elem, _ := doc.CreateElement("elem")
	pi, _ := doc.CreateProcessingInstruction("target", "data")

	doc.AppendChild(root)
	root.AppendChild(childOfRoot)
	doc.AppendChild(pi)

	replacement, err := doc.ReplaceChild(elem, root)
	if err != nil {
		t.Error("unexpected error")
	}

	// The replacement node must equal the replaced node.
	if replacement != root {
		t.Error("replacement should equal elem")
	}

	if doc.GetDocumentElement() != elem {
		t.Error("incorrect document element")
	}

	// Try to replace the processing instruction with another element. Should fail
	// because we cannot have two document elements.
	if _, err = doc.ReplaceChild(root, pi); err == nil {
		t.Error("expected an error but got none")
	}

	// Try to replace the current root node (elem) with childOfRoot.
	if _, err = doc.ReplaceChild(childOfRoot, elem); err != nil {
		t.Error("unexpected error")
	}

	if doc.GetDocumentElement() != childOfRoot {
		t.Error("incorrect document element")
	}

	if len(doc.GetChildNodes()) != 2 {
		t.Error("expected 2 children")
	}
}
