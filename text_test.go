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
	if txt.GetLastChild() != nil {
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
	if txt.GetNamespaceURI() != "" {
		t.Error("namespace URI should be an empty string")
	}
	if txt.GetNamespacePrefix() != "" {
		t.Error("namespace prefix should be an empty string")
	}
	if _, err := txt.ReplaceChild(nil, nil); err == nil {
		t.Error("replacing child should always return an error")
	}
	if _, err := txt.RemoveChild(nil); err == nil {
		t.Error("removing child should always return an error")
	}
	if _, err := txt.InsertBefore(nil, nil); err == nil {
		t.Error("inserting a child should always return an error")
	}
	if txt.LookupPrefix("anything") != "" {
		t.Error("LookupPrefix should always return an empty string")
	}
	if _, found := txt.LookupNamespaceURI("asd"); found {
		t.Error("LookupNamespaceURI should always return an empty string and false")
	}
	thetext := "some text content"
	txt.SetTextContent(thetext)
	if txt.GetTextContent() != thetext {
		t.Errorf("expected '%v', got '%v'", thetext, txt.GetTextContent())
	}

	txt.SetText("this is a sample string, longer than 30 characters.")
	s := fmt.Sprintf("%v", txt)
	expected := "TEXT_NODE: 'this is a sample string, longe [...]'"
	if s != expected {
		t.Errorf("expected '%v', got '%v'", expected, s)
	}
}

func TestTextPrevNextSibling(t *testing.T) {
	doc := NewDocument()
	root, _ := doc.CreateElement("root")
	child1, _ := doc.CreateElement("child1")
	text := doc.CreateText("some text")
	child3, _ := doc.CreateElement("child3")

	doc.AppendChild(root)
	root.AppendChild(child1)
	root.AppendChild(text)
	root.AppendChild(child3)

	node := text.GetPreviousSibling()
	if node != child1 {
		t.Errorf("expected 'child1' as previous sibling, got '%v'", node)
	}
	node = text.GetNextSibling()
	if node != child3 {
		t.Errorf("expected 'child3' as next sibling, got '%v'", node)
	}
}

func TestTextIsIgnorableWhitespace(t *testing.T) {
	var tests = []struct {
		t        string
		expected bool
	}{
		{"\n\r\n\r\t\t", true},
		{"\n\ra\n\r\t\t", false},
		{"\r\n\r\n   \r    ", true},
		{"         ", true},
		{"         \r", true},
		{"     x            ", false},
	}

	doc := NewDocument()
	for _, test := range tests {
		text := doc.CreateText(test.t)
		if text.IsElementContentWhitespace() != test.expected {
			t.Errorf("expected '%v'", test.expected)
		}
	}
}
