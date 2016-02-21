package dom

import (
	"encoding/xml"
	"fmt"
	"io"
)

type Builder struct {
	reader  io.Reader
	doc     Document
	decoder *xml.Decoder
}

func NewBuilder(reader io.Reader) *Builder {
	b := &Builder{}
	b.reader = reader
	b.doc = NewDocument()
	b.decoder = xml.NewDecoder(b.reader)
	return b
}

func (b *Builder) CreateDocument() (Document, error) {
	var curNode Node = b.doc
	_ = curNode
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
		case xml.StartElement:
			fmt.Printf("Start elem: %s\n", typ.Name)
		}
	}
	return b.doc, nil
}

func (b *Builder) derp(node Node) Node {
	return node
}
