package dom

import (
	"testing"
)

func TestElementNodeName(t *testing.T) {
	root := newElement()
	root.SetTagName("rewt")

	if root.GetNodeName() != "rewt" {
		t.Errorf("node name is expected to be the same as the tag name, which is '%v'", root.GetTagName())
	}
}

func TestElementNodeType(t *testing.T) {
	root := newElement()
	if root.GetNodeType() != ElementNode {
		t.Errorf("elements are supposed to be of type '%v'", ElementNode)
	}
}

func TestElementNodeValue(t *testing.T) {
	root := newElement()
	if root.GetNodeValue() != "" {
		t.Errorf("node value of elements are not applicable and should therefore be empty")
	}
}

func TestElementLocalName(t *testing.T) {
	// TODO: check out the local name
}

func TestElementNodeList(t *testing.T) {
	root := newElement()
	if len(root.GetChildNodes()) != 0 {
		t.Errorf("uninitialized node list should be zero length")
	}

	// add some children
	for i := 0; i < 10; i++ {
		e := newElement()
		e.SetTagName("element" + string(i))
		root.AppendChild(e)
	}

	if len(root.GetChildNodes()) != 10 {
		t.Errorf("node list length should be 10, but was ", len(root.GetChildNodes()))
	}

}

func TestElementOwnerDocument(t *testing.T) {
	doc := NewDocument()
	elem, err := doc.CreateElement("root")
	if err != nil {
		t.Errorf("no error expected!!")
	}
	if elem.GetOwnerDocument() != doc {
		t.Errorf("incorrect owner document")
	}
}

func TestElementAppendChild(t *testing.T) {
	root := newElement()
	root.SetTagName("root")

	child := newElement()
	child.SetTagName("child")

	if len(root.GetChildNodes()) != 0 {
		t.Errorf("length of node list should be 0, but was %v", len(root.GetChildNodes()))
	}

	err := root.AppendChild(child)
	if err != nil {
		t.Errorf("did not expect error at this point: '%v'", err)
	}

	if child.GetParentNode() != root {
		t.Logf("child.parent = %v, root = %v", child.GetParentNode(), root)
		t.Errorf("parent node of 'child' should be 'root'")
	}

	if len(root.GetChildNodes()) != 1 {
		t.Errorf("length of node list should be 1, but was %v", len(root.GetChildNodes()))
	}

	if root.GetChildNodes()[0] != child {
		t.Errorf("first element is expected to be 'child', but was %v", root.GetChildNodes()[0])
	}

	err = root.AppendChild(root)
	if err == nil {
		t.Errorf("expected a hierarchy error here")
	}
}

func TestElementHasChildNodes(t *testing.T) {
	root := newElement()
	if root.HasChildNodes() {
		t.Errorf("expected no child nodes")
	}

	// Adding root to itself should result in nothing
	root.AppendChild(root)
	if root.HasChildNodes() {
		t.Errorf("expected no child nodes")
	}

	child := newElement()
	root.AppendChild(child)
	if !root.HasChildNodes() {
		t.Errorf("expected child nodes")
	}
}
