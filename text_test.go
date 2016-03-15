package dom

import (
	"fmt"
	"testing"
)

func TestTextGetters(t *testing.T) {
	doc := NewDocument()
	parent, _ := doc.CreateElement("element")
	txt := doc.CreateText("sample")
	doc.AppendChild(parent)
	parent.AppendChild(txt)

	if txt.GetNodeName() != "#text" {
		t.Error("expected '#text'")
	}
	if txt.GetNodeType() != TextNode {
		t.Error("expected TextNode")
	}
	if txt.GetNodeValue() != "sample" {
		t.Errorf("expected 'sample', got '%v'", txt.GetNodeValue())
	}
	if txt.GetLocalName() != "" {
		t.Error("local name should always be an empty string")
	}
	if txt.GetParentNode() != parent {
		t.Error("incorrect parent node")
	}
	if err := txt.AppendChild(doc.CreateText("meh")); err == nil {
		t.Error("expected error, but got none")
	}
	if txt.GetFirstChild() != nil {
		t.Error("text nodes cannot have children")
	}
	if txt.HasChildNodes() {
		t.Error("text nodes cannot have children")
	}
	if txt.GetAttributes() != nil {
		t.Error("text nodes cannot have attributes")
	}
	if txt.GetOwnerDocument() != doc {
		t.Error("incorrect owner document")
	}
	if txt.GetNamespacePrefix() != "" {
		t.Error("namespace prefix should be an empty string")
	}
	if txt.LookupPrefix("anything") != "" {
		t.Error("LookupPrefix should always return an empty string")
	}
	if txt.LookupNamespaceURI("asd") != "" {
		t.Error("LookupNamespaceURI should always return an empty string")
	}

	txt.SetText("this is a sample string, longer than 30 characters.")
	s := fmt.Sprintf("%v", txt)
	expected := "TEXT_NODE: 'this is a sample string, longe [...]'"
	if s != expected {
		t.Errorf("expected '%v', got '%v'", expected, s)
	}
}
