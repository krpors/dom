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
func TestParserParse(t *testing.T) {
	reader := strings.NewReader(exampleDoc1)
	builder := NewParser(reader)
	doc, err := builder.Parse()
	if err != nil {
		t.Errorf("unexpected error after building document from string: %v", err)
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

	// Get first person elemen, check attributes.
	elements := doc.GetElementsByTagName("person")
	if len(elements) <= 0 {
		t.Error("expected at least 1 element")
		t.FailNow()
	}

	attrVal := elements[0].GetAttribute("name")
	if attrVal != "Foo" {
		t.Errorf("expected 'Foo', got '%v'", attrVal)
	}
	attrVal = elements[0].GetAttribute("lastname")
	if attrVal != "Quux" {
		t.Errorf("expected 'Quux', got ''%v'", attrVal)
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
func TestParserParseWhitespaces(t *testing.T) {
	reader := strings.NewReader(exampleDoc2)
	builder := NewParser(reader)
	doc, err := builder.Parse()
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

func TestParserParseContentInProlog(t *testing.T) {
	reader := strings.NewReader(exampleErrDoc1)
	builder := NewParser(reader)
	_, err := builder.Parse()
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

func TestParserParseTrailingChars(t *testing.T) {
	reader := strings.NewReader(exampleErrDoc2)
	builder := NewParser(reader)
	_, err := builder.Parse()
	if err == nil {
		t.Errorf("expected error due to trailing character data after root node")
	}
}

//=============================================================================

// exampleDoc3 contains elaborate namespaces for testing that.
var exampleDoc3 = `<?xml version="1.0" encoding="UTF-8"?>
<!-- rootNode has the default namespace with some declarations of namespaces -->
<rootNode xmlns:pfx="urn:ns:pfx" xmlns="urn:rootnode">
    <!-- This element's namespace should match "urn:ns:pfx:childelement": -->
    <pfx:childElement xmlns:pfx="urn:ns:pfx:childelement">
        <!-- This element's namespace should match "urn:ns:pfx:sameprefix": -->
        <pfx:samePrefix xmlns:pfx="urn:ns:pfx:sameprefix"/>
    </pfx:childElement>
    <pfx:childElement>
        <!-- This element's namespace should match "urn:ns:pfx" -->
    </pfx:childElement>
</rootNode>`

// Tests whether embedding the same prefixes (pfx) within a document results
// in the correct associated namespaces.
func TestParserParseNamespaces(t *testing.T) {
	reader := strings.NewReader(exampleDoc3)
	builder := NewParser(reader)
	doc, err := builder.Parse()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	docelem := doc.GetDocumentElement()
	if docelem.GetNamespaceURI() != "urn:rootnode" {
		t.Errorf("namespace URI of root node should be 'urn:rootnode', but was '%v'", docelem.GetNamespaceURI())
	}

	childElement1 := docelem.GetChildNodes()[3]
	if childElement1.GetNodeName() != "pfx:childElement" {
		t.Errorf("expected 'childElement', got '%v'", childElement1.GetNodeName())
	}
	if childElement1.GetNamespaceURI() != "urn:ns:pfx:childelement" {
		t.Errorf("expected 'urn:ns:pfx:childelement', got '%v'", childElement1.GetNamespaceURI())
	}

	samePrefix := childElement1.GetChildNodes()[3]
	if samePrefix.GetNodeName() != "pfx:samePrefix" {
		t.Errorf("expected 'samePrefix', got '%v'", samePrefix.GetNodeName())
	}

	if samePrefix.GetNamespaceURI() != "urn:ns:pfx:sameprefix" {
		t.Errorf("expected 'urn:ns:pfx:sameprefix', got '%v'", samePrefix.GetNamespaceURI())
	}

	childElement2 := docelem.GetChildNodes()[5]
	if childElement2.GetNodeName() != "childElement" {
		t.Errorf("expected 'childElement', got '%v'", childElement1.GetNodeName())
	}

	if childElement2.GetNamespaceURI() != "urn:ns:pfx" {
		t.Errorf("expected 'urn:ns:pfx', got '%v'", childElement2.GetNamespaceURI())
	}

	// Try to find some elements using GetElementsByTagName[NS].
	gebtn := doc.GetElementsByTagName("pfx:childElement")
	if len(gebtn) != 1 {
		t.Errorf("expected 1, got %d", len(gebtn))
	}
	meh := gebtn[0].GetElementsByTagNameNS("urn:ns:pfx:sameprefix", "samePrefix")
	if len(meh) != 1 {
		t.Errorf("expected 1, got %d", len(meh))
	}
	gebtn = docelem.GetElementsByTagNameNS("urn:ns:pfx:childelement", "childElement")
	if len(gebtn) != 1 {
		t.Errorf("expected 1, got %d", len(gebtn))
	}
	gebtn = docelem.GetElementsByTagNameNS("urn:ns:pfx", "childElement")
	if len(gebtn) != 1 {
		t.Errorf("expected 1, got %d", len(gebtn))
	}
}

var exampleDoc4 = `<?xml version="1.0" encoding="UTF-8"?>
<soap:Envelope xmlns:soap="http://www.w3.org/2003/05/soap-envelope">
  <soap:Header xmlns:soap="http://www.w3.org/2003/05/soap-envelope" xmlns:superfluous="http://www.w3.org/2003/05/soap-envelope">
  </soap:Header>
  <soap:Body xmlns:soap="http://www.w3.org/2003/05/soap-envelope">
    <m:GetStockPrice xmlns:m="http://www.example.org/stock/Surya">
      <m:StockName xmlns:m="http://www.example.org/stock/Surya" identifier="cruft">IBM and more!</m:StockName>
	  <soap:Extension>Wow, what's happening!</soap:Extension>
    </m:GetStockPrice>
  </soap:Body>
</soap:Envelope>`

func TestParserWut(t *testing.T) {
	reader := strings.NewReader(exampleDoc4)
	parser := NewParser(reader)
	parser.Configuration.ElementContentWhitespace = false
	parser.Configuration.Comments = false
	doc, err := parser.Parse()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	doc.NormalizeDocument()
	// TODO: verify?
}
