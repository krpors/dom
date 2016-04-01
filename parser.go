package dom

import (
	"encoding/xml"
	"fmt"
	"io"
	"strings"
)

// Parser is the entrypoint of the dom package to parse an XML tree from the given
// reader, into the doc attribute, using the decoder from encoding/xml.
type Parser struct {
	reader io.Reader // Reader containing the XML document.

	Configuration Configuration
}

// NewParser constructs a new Parser using the given reader. The reader is expected
// to contain the (...valid) XML tree. Namespace awareness will be set to true per default.
// Parser configuration will be set to a default one.
func NewParser(reader io.Reader) *Parser {
	b := &Parser{}
	b.reader = reader
	b.Configuration = NewConfiguration()
	return b
}

// Parse parses an XML Document contained within the reader attribute of the current Parser.
// A Document will be returned and a nil error if the parsing succeeded.
func (b *Parser) Parse() (Document, error) {
	doc := NewDocument()
	decoder := xml.NewDecoder(b.reader)
	var curNode = Node(doc)

	for {
		token, err := decoder.Token()
		if err == io.EOF {
			// End of file, processed okay
			return doc, nil
		}
		if err != nil {
			// Other error, return that.
			return nil, err
		}

		switch typ := token.(type) {
		case xml.Comment:
			// Skip comments?
			if !b.Configuration.Comments {
				continue
			}

			cmt, err := doc.CreateComment(string(typ))
			if err != nil {
				return nil, err
			}
			if err = curNode.AppendChild(cmt); err != nil {
				return nil, err
			}
		case xml.ProcInst:
			// Note: the Go default decoder regards the XML declaration as a processing
			// instruction, even though it is not. Therefore, we handle this edge case
			// to NOT include this as a valid child node.
			if strings.ToLower(typ.Target) != "xml" {
				pi, err := doc.CreateProcessingInstruction(typ.Target, string(typ.Inst))
				if err != nil {
					return nil, err
				}
				if err = doc.AppendChild(pi); err != nil {
					return nil, err
				}
			}
		case xml.StartElement:
			//  FIXME: The default encoding/xml.Decoder does fuck all about prefixes.
			// That's not all: https://github.com/golang/go/issues/11735
			// Therefore, we don't/can't set any kind of prefix ourselves. The
			// NormalizeDocument() or Normalize() can be used to "fix" the namespaces.
			namespace := ""
			// Are we parsing namespaces? e.g. namespace awareness?
			if b.Configuration.NamespaceDeclarations {
				namespace = typ.Name.Space
			}

			elem, err := doc.CreateElementNS(namespace, typ.Name.Local)
			if err != nil {
				return nil, err
			}

			// Iterate over the element's attributes.
			for _, a := range typ.Attr {
				namespace = ""
				// Are we parsing namespaces?
				if b.Configuration.NamespaceDeclarations {
					namespace = a.Name.Space
				}

				attrName := a.Name.Local

				// OK: there is some weirdness going in with parsing attributes.
				// For example, these have different outputs:
				//
				//	xmlns:bla="what" : namespace = "xmlns", localname = "bla", value = "what".
				//  bla="what" : namespace = "" (unless bound), localname = "bla", value = "what"
				//  ns0:bla="what" : namespace = [whatever ns0 is bound to], localname = "bla"
				//
				// Therefore, we handle the xmlns edge case like this, so we can assign prefixes
				// to elements. Looks extremely hacky, and it is. I wish it could behave more SAX like.

				// fmt.Printf("%s, %s, %s\n", a.Name.Space, a.Name.Local, a.Value)

				if strings.HasPrefix(a.Name.Space, "xmlns") {
					if a.Value == typ.Name.Space {
						// At this point, there is a match between a namespace declaration (+ prefix) append
						// the StartElement we're in. We're gonna give our Element a namespace prefix.
						elem.setTagName(a.Name.Local + ":" + typ.Name.Local)
					} else {
						elem.SetAttribute("xmlns:"+a.Name.Local, a.Value)
					}
					continue // Skip it!! No need for namespace declarations anyway. They can be fixed by normalizing.
				}

				// Add all other (normal) attributes.
				attr, err := doc.CreateAttributeNS(namespace, attrName)
				if err != nil {
					return nil, err
				}
				attr.SetValue(a.Value)
				elem.GetAttributes().SetNamedItem(attr)
			}

			if err = curNode.AppendChild(elem); err != nil {
				return nil, err
			}
			curNode = elem
		case xml.EndElement:
			curNode = curNode.GetParentNode()
		case xml.CharData:
			// If there is no document element yet, and the character data is found which is NOT whitespace,
			// generate an error. No character data allowed before document element, but whitespaces
			// are okay to parse. Don't add it as a child element though.
			if doc.GetDocumentElement() == nil {
				if strings.TrimSpace(string(typ)) != "" {
					return nil, fmt.Errorf("%v: content is not allowed in prolog", ErrorHierarchyRequest)
				}
				// We got whitespace. Don't add it as a child, merely continue the next token
				// parsing in the stream.
				continue
			}
			// Likewise, character data may not occur after the document element in the trailing
			// section, so check that as well. The Go decoder doesn't care so we handle this edge
			// case as well.
			if curNode == doc {
				if strings.TrimSpace(string(typ)) != "" {
					// We cannot append text/chardata to the document itself.
					return nil, fmt.Errorf("%v: content is not allowed in trailing section", ErrorHierarchyRequest)
				}
				// We got whitespace. Don't add it as a child, merely continue the next token
				// parsing in the stream. Same behaviour as above.
				continue
			}

			text := doc.CreateText(string(typ))
			// Should we ignore ignorable whitespaces, and the text content is whitespace?
			if !b.Configuration.ElementContentWhitespace && text.IsElementContentWhitespace() {
				continue
			}

			// In all other cases, create a text node and add it to the current node as a child.
			if err := curNode.AppendChild(text); err != nil {
				return nil, err
			}
		}
	}
}
