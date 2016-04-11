package dom

import (
	"testing"
)

func TestProcessingInstructionGetters(t *testing.T) {
	doc := NewDocument()

	pi, _ := doc.CreateProcessingInstruction("procinsttarget", "procinstdata")
	pi.setParentNode(doc)

	bogusElement, _ := doc.CreateElement("anything")
	if err := pi.AppendChild(bogusElement); err == nil {
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
	if pi.GetLastChild() != nil {
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
	if pi.GetNamespacePrefix() != "" {
		t.Error("namespace prefix should be an empty string")
	}
	if _, err := pi.RemoveChild(nil); err == nil {
		t.Error("expected error, but got none")
	}
	if _, err := pi.InsertBefore(nil, nil); err == nil {
		t.Error("expected error, but got none")
	}
	if _, err := pi.ReplaceChild(nil, nil); err == nil {
		t.Error("expected error, but got none")
	}
	if pi.HasAttributes() {
		t.Error("processing instructions cannot have children")
	}
}
