package dom

import (
	"strings"
	"testing"
)

// exampleErrDoc1 contains an invalid XML document, because character data exists
// between the XML declaration and the document element.
var exampleErrDoc1 = `<?xml version="1.0" encoding="UTF-8"?>


invalid content.

<directory>
	Character data can exist here.
</directory>`

// exampleDoc contains an XML valid document.
var exampleDoc = `<?xml version="1.0" encoding="UTF-8"?>



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

func TestBuilderCreateDocument(t *testing.T) {
	reader := strings.NewReader(exampleDoc)
	builder := NewBuilder(reader)
	doc, err := builder.CreateDocument()
	if err != nil {
		t.Errorf("unexpected error after building document from string: '%v'", err)
		t.FailNow()
	}

	if doc.GetFirstChild().GetNodeName() != "directory" {
		t.Errorf("expecting root node 'directory', but was '%s'", doc.GetFirstChild().GetNodeName())
	}

	// Try navigating directly to the first comment:
	cmt := doc.
		GetFirstChild().    // <directory>
		GetChildNodes()[1]. // [0] = text node, [1] = <people>
		GetChildNodes()[5]  // [0] = text, [1] = <person>, [2] = text, [3] = <person>, [4] = text, [5] = comment node

	if cmt.GetNodeType() != CommentNode {
		t.Errorf("expecting a comment node, but was %v", cmt.GetNodeType())
	}
}

func TestBuilderCreateDocumentInvalidContent(t *testing.T) {
	reader := strings.NewReader(exampleErrDoc1)
	builder := NewBuilder(reader)
	_, err := builder.CreateDocument()
	if err == nil {
		t.Errorf("expected error after building document from string, but got none")
	}
}
