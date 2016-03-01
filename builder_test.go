package dom

import (
	"strings"
	"testing"
)

// exampleDoc1 contains an XML valid document.
var exampleDoc1 = `<?xml version="1.0" encoding="UTF-8"?>
<!-- Comment may occur here. -->
<directory>
	<people>
		<person name="Foo" lastname="Quux">
			<date>2016-01-01</date>
		</person>
		<person name="Alice" lastname="Bob">
			<date>2016-09-03</date>
		</person>
		<!-- this is a comment -->
		<!-- empty element follows -->
		<person/>
	</people>
	<ns:cruft xmlns:ns="http://example.org/xmlns/uri">
		<ns:other>Character data.</ns:other>
		<ns:balls ns:derp="woot">More chardata</ns:balls>
	</ns:cruft>
	<Grøups>asd</Grøups>
</directory>`

// Tests a completely valid document and checks whether everything is in place.
func TestBuilderCreateDocument(t *testing.T) {
	reader := strings.NewReader(exampleDoc1)
	builder := NewBuilder(reader)
	doc, err := builder.CreateDocument()
	if err != nil {
		t.Errorf("unexpected error after building document from string: '%v'", err)
		t.FailNow()
	}

	if doc.GetFirstChild().GetNodeType() != CommentNode {
		t.Errorf("expecting comment node as first child, but was '%v'", doc.GetFirstChild())
	}

	if doc.GetChildNodes()[1].GetNodeName() != "directory" {
		t.Errorf("expecting root node 'directory', but was '%s'", doc.GetFirstChild().GetNodeName())
		t.FailNow()
	}

	// Try navigating directly to the first comment:
	cmt := doc.
		GetDocumentElement(). // <directory>
		GetChildNodes()[1].   // [0] = text node, [1] = <people>
		GetChildNodes()[5]    // [0] = text, [1] = <person>, [2] = text, [3] = <person>, [4] = text, [5] = comment node

	if cmt.GetNodeType() != CommentNode {
		t.Errorf("expecting a comment node, but was %v", cmt.GetNodeType())
	}
}

//=============================================================================

// exampleDoc2 contains a valid XML document with leading and trailing whitespaces
// before and after the root node (the document element).
var exampleDoc2 = `<?xml version="1.0" encoding="UTF-8"?>


<whitespaces/>


`

// Tests whether whitespaces before and after the document element don't generate
// an error.
func TestBuilderCreateDocumentWhitespaces(t *testing.T) {
	reader := strings.NewReader(exampleDoc2)
	builder := NewBuilder(reader)
	doc, err := builder.CreateDocument()
	if err != nil {
		t.Errorf("unexpected error: '%v'", err)
		t.FailNow()
	}

	// Check the correct amount of children. Make sure no three children are added,
	// just the document element (so no text, docelem, text).
	if len(doc.GetChildNodes()) != 1 {
		t.Errorf("expecting one child node but got %v", len(doc.GetChildNodes()))
	}

	if doc.GetFirstChild().GetNodeName() != "whitespaces" {
		t.Errorf("expecting 'whitespaces' as nodename")
	}
}

//=============================================================================

// exampleErrDoc1 contains an invalid XML document, because character data exists
// between the XML declaration and the document element.
var exampleErrDoc1 = `<?xml version="1.0" encoding="UTF-8"?>


invalid content.

<directory>
	Character data can exist here.
</directory>`

func TestBuilderCreateDocumentContentInProlog(t *testing.T) {
	reader := strings.NewReader(exampleErrDoc1)
	builder := NewBuilder(reader)
	_, err := builder.CreateDocument()
	if err == nil {
		t.Errorf("expected error after building document from string, but got none")
	}
}

//=============================================================================

// exampleErroDoc2 contains trailing character data after the document, which
// should result in an error.
var exampleErrDoc2 = `<?xml version="1.0" encoding="UTF-8"?>
<stuff>
</stuff>
chardata content in trailing section`

func TestBuilderCreateDocumentTrailingChars(t *testing.T) {
	reader := strings.NewReader(exampleErrDoc2)
	builder := NewBuilder(reader)
	_, err := builder.CreateDocument()
	if err == nil {
		t.Errorf("expected error due to trailing character data after root node")
	}
}

//=============================================================================

// exampleDoc3 contains elaborate namespaces for testing that.
var exampleDoc3 = `<?xml version="1.0" encoding="UTF-8"?>
<root xmlns="http://example.org/rooturi">
	<pfx:child1 xmlns:pfx="urn:narf:zoit">
		<pfx:subchild>This element should be of namespace urn:narf:zoit.</pfx:subchild>
	</pfx:child1>
	<same_root_namespace/>
</root>
`

func TestBuilderCreateDocumentNamespaces(t *testing.T) {
	reader := strings.NewReader(exampleDoc3)
	builder := NewBuilder(reader)
	doc, err := builder.CreateDocument()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	docelem := doc.GetDocumentElement()
	if docelem.GetNamespaceURI() != "http://example.org/rooturi" {
		t.Errorf("namespace URI of root node was '%v', expected 'http://example.org/rooturi'", docelem.GetNamespaceURI())
	}

	child1 := docelem.GetChildNodes()[1].(Element)
	if child1.GetNamespaceURI() != "urn:narf:zoit" {
		t.Errorf("namespace URI of first element should be 'urn:narf:zoit', but was '%v'", child1.GetNamespaceURI())
	}
	if child1.GetNamespacePrefix() != "pfx" {
		t.Errorf("namespace prefix should be 'pfx', but was '%v'", child1.GetNamespacePrefix())
	}
}
