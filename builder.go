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
		// TODO: attributes
		case xml.Comment:
			cmt, err := b.doc.CreateComment(string(typ))
			if err != nil {
				return nil, err
			}
			curNode.AppendChild(cmt)
		case xml.ProcInst:
			// Note: the Go default decoder regards the XML declaration as a processing
			// instruction, even though it is not. Therefore, we handle this edge case
			// to NOT include this as a valid child node.
			if strings.ToLower(typ.Target) != "xml" {
				pi, err := b.doc.CreateProcessingInstruction(typ.Target, string(typ.Inst))
				if err != nil {
					return nil, err
				}
				b.doc.AppendChild(pi)
			}
		case xml.StartElement:
			elem, err := b.doc.CreateElementNS(typ.Name.Space, typ.Name.Local)
			if err != nil {
				return nil, err
			}
			curNode.AppendChild(elem)
			curNode = elem
		case xml.EndElement:
			curNode = curNode.GetParentNode()
		case xml.CharData:
			// FIXME: character data is still read by the xml.Decoder even AFTER the document element.
			// This isn't good. Return an error then.
			text := b.doc.CreateTextNode(string(typ))
			curNode.AppendChild(text)
		}
	}
}

func (b *Builder) attrsToBleh(a []xml.Attr) Attr {
	return nil
}
