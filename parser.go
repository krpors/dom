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

	NamespaceAware          bool // Enable namespace awareness. Default: true.
	ReadIgnorableWhitespace bool // Enable ignoring of ignorable whitespace. Default: true.
}

// NewParser constructs a new Parser using the given reader. The reader is expected
// to contain the (...valid) XML tree. Namespace awareness will be set to true per default.
func NewParser(reader io.Reader) *Parser {
	b := &Parser{}
	b.reader = reader
	b.NamespaceAware = true
	b.ReadIgnorableWhitespace = true
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
			// Therefore, we don't set any kind of prefix ourselves.
			namespace := ""
			if b.NamespaceAware {
				namespace = typ.Name.Space
			}
			elem, err := doc.CreateElementNS(namespace, typ.Name.Local)
			if err != nil {
				return nil, err
			}

			for _, a := range typ.Attr {
				namespace = ""
				if b.NamespaceAware {
					namespace = a.Name.Space
				}
				attr, err := doc.CreateAttributeNS(namespace, a.Name.Local)
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
			if !b.ReadIgnorableWhitespace && text.IsElementContentWhitespace() {
				continue
			}

			// In all other cases, create a text node and add it to the current node as a child.
			if err := curNode.AppendChild(text); err != nil {
				return nil, err
			}
		}
	}
}