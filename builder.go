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
	b.decoder = xml.NewDecoder(b.reader)
	return b
}

func (b *Builder) PrintTree(w io.Writer) {

	var xtree func(n Node, padding string)
	xtree = func(n Node, padding string) {
		w.Write([]byte(fmt.Sprintf("%s%s\n", padding, n)))
		for _, node := range n.NodeList() {
			xtree(node, padding+"  ")
		}
	}

	xtree(b.doc, "")
}

func (b *Builder) CreateDocument() (Document, error) {
	b.doc = NewDocument()
	var curNode Node = b.doc

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
		case xml.ProcInst:
			//fmt.Println(string(typ.Target))// TODO: processing instruction type.
			//fmt.Println(string(typ.Inst))
		case xml.StartElement:
			elem, err := b.doc.CreateElementNS(typ.Name.Space, typ.Name.Local)
			if err != nil {
				return nil, err
			}
			curNode.AppendChild(elem)
			curNode = elem
		case xml.EndElement:
			curNode = curNode.ParentNode()
		case xml.CharData:
			text := b.doc.CreateTextNode(string(typ))
			curNode.AppendChild(text)
		}
	}

	return b.doc, nil
}

func (b *Builder) attrsToBleh(a []xml.Attr) Attr {
	return nil
}
