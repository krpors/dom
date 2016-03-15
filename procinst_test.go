package dom

import (
	"testing"
)

func TestProcessingInstructionGetters(t *testing.T) {
	doc := NewDocument()

	pi := newProcInst(doc)
	pi.setData("procinstdata")
	pi.setTarget("procinsttarget")
	pi.setNamespaceURI("http://example.org/unused/and/a/no/op")
	pi.setParentNode(doc)

	if err := pi.AppendChild(newElement(doc)); err == nil {
		t.Error("expected error at this point")
	}
	if pi.GetOwnerDocument() != doc {
		t.Error("owner document invalid")
	}
	if pi.GetParentNode() != doc {
		t.Error("invalid parent node (should be document)")
	}
	if pi.GetFirstChild() != nil {
		t.Error("processing instructions should not have children")
	}
	if len(pi.GetChildNodes()) != 0 {
		t.Error("child nodes length should be 0")
	}
	if pi.HasChildNodes() {
		t.Error("HasChildNodes returns true, should be false")
	}
	if pi.GetAttributes() != nil {
		t.Error("attributes should be nil")
	}
	if pi.GetTarget() != "procinsttarget" {
		t.Errorf("target should be 'procinsttarget', but was '%v'", pi.GetTarget())
	}
	if pi.GetNodeName() != "procinsttarget" {
		t.Errorf("nodename should be 'procinsttarget', but was '%v'", pi.GetNodeName())
	}
	if pi.GetData() != "procinstdata" {
		t.Errorf("data should be 'procinstdata', but was '%v'", pi.GetData())
	}
	if pi.GetNodeValue() != "procinstdata" {
		t.Errorf("node value should be 'procinstdata', but was '%v'", pi.GetNodeValue())
	}
	if pi.GetNamespaceURI() != "" {
		t.Error("namespace uri should be an empty string at all times")
	}
	if pi.GetNodeType() != ProcessingInstructionNode {
		t.Errorf("node type should be '%v'", pi.GetNodeType())
	}
}
