package dom

import (
	"testing"
)

func TestElementGetters(t *testing.T) {
	root := newElement()
	root.SetTagName("pfx:rewt")

	if root.GetNodeName() != "pfx:rewt" {
		t.Errorf("node name is expected to be 'pfx:rewt', but was '%v'", root.GetNodeName())
	}
	if root.GetNodeType() != ElementNode {
		t.Errorf("elements are supposed to be of type '%v'", ElementNode)
	}
	if root.GetNodeValue() != "" {
		t.Errorf("node value of elements are not applicable and should therefore be empty")
	}
	if len(root.GetChildNodes()) != 0 {
		t.Errorf("initialized node list should be zero length")
	}
	if root.GetNodeName() != root.GetTagName() {
		t.Errorf("GetNodeName() should equal GetTagName(): %v != %v", root.GetNodeName(), root.GetTagName())
	}
	if root.GetLocalName() != "rewt" {
		t.Errorf("local name should be 'rewt', but was '%v'", root.GetLocalName())
	}
	if root.GetNamespacePrefix() != "pfx" {
		t.Errorf("namespace prefix should be 'pfx', but was '%v'", root.GetNamespacePrefix())
	}
	// add some children
	for i := 0; i < 10; i++ {
		e := newElement()
		e.SetTagName("element" + string(i))
		root.AppendChild(e)
	}
	if len(root.GetChildNodes()) != 10 {
		t.Errorf("node list length should be 10, but was %v", len(root.GetChildNodes()))
	}

	// element without prefix:
	root = newElement()
	root.SetTagName("rewt")
	if root.GetLocalName() != "rewt" {
		t.Errorf("local name should be 'rewt', was '%v'", root.GetLocalName())
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

	if root.GetFirstChild() != nil {
		t.Error("first child should be nil")
	}

	err := root.AppendChild(child)
	if err != nil {
		t.Errorf("did not expect error at this point: '%v'", err)
	}

	if root.GetFirstChild() != child {
		t.Error("invalid first child")
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
