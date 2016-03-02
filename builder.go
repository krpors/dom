package dom

import (
	"encoding/xml"
	"fmt"
	"io"
	"strings"
)

// Builder is the entrypoint of the dom package to parse an XML tree from the given
// reader, into the doc attribute, using the decoder from encoding/xml.
type Builder struct {
	reader  io.Reader
	doc     Document
	decoder *xml.Decoder
}

// NewBuilder constructs a new Builder using the given reader. The reader is expected
// to contain the (...valid) XML tree.
func NewBuilder(reader io.Reader) *Builder {
	b := &Builder{}
	b.reader = reader
	b.decoder = xml.NewDecoder(b.reader)
	return b
}

// PrintTree is a utility function to print the parsed document to an internal
// representation of the complete hierarchy. Needs work.
func (b *Builder) PrintTree(w io.Writer) {

	var xtree func(n Node, padding string)
	xtree = func(n Node, padding string) {
		w.Write([]byte(fmt.Sprintf("%s%v\n", padding, n)))
		for _, node := range n.GetChildNodes() {
			xtree(node, padding+"  ")
		}
	}

	xtree(b.doc, "")
}

// CreateDocument creates a Document object using the constructed decoder.
// Will return the Document if everything went a-okay, or a non-nil error
// if something has failed during the parsing of the tokens.
func (b *Builder) CreateDocument() (Document, error) {
	b.doc = NewDocument()
	var curNode = Node(b.doc)

	for {
		token, err := b.decoder.Token()
		if err == io.EOF {
			// End of file, processed okay
			return b.doc, nil
		}
		if err != nil {
			// Other error, return that.
			return nil, err
		}

		switch typ := token.(type) {
		case xml.Attr:
			// attr, err := b.doc.CreateAttribute()
		case xml.Comment:
			cmt, err := b.doc.CreateComment(string(typ))
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
				pi, err := b.doc.CreateProcessingInstruction(typ.Target, string(typ.Inst))
				if err != nil {
					return nil, err
				}
				if err = b.doc.AppendChild(pi); err != nil {
					return nil, err
				}
			}
		case xml.StartElement:
			//  FIXME: The default encoding/xml.Decoder does fuck all about prefixes.
			// That's not all: https://github.com/golang/go/issues/11735
			// Therefore, we don't set any kind of prefix ourselves.
			elem, err := b.doc.CreateElementNS(typ.Name.Space, typ.Name.Local)
			if err != nil {
				return nil, err
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
			if b.doc.GetDocumentElement() == nil {
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
			if curNode == b.doc {
				if strings.TrimSpace(string(typ)) != "" {
					// We cannot append text/chardata to the document itself.
					return nil, fmt.Errorf("%v: content is not allowed in trailing section", ErrorHierarchyRequest)
				}
				// We got whitespace. Don't add it as a child, merely continue the next token
				// parsing in the stream. Same behaviour as above.
				continue
			}
			// In all other cases, create a text node and add it to the current node as a child.
			text := b.doc.CreateText(string(typ))
			if err := curNode.AppendChild(text); err != nil {
				return nil, err
			}
		}
	}
}

func (b *Builder) attrsToBleh(a []xml.Attr) Attr {
	return nil
}
