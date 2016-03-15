package dom

import (
	"testing"
)

// Tests the getters of a Comment node.
func TestCommentGetters(t *testing.T) {
	doc := NewDocument()
	cmt, err := doc.CreateComment("comment string")
	if err != nil {
		t.Error("unexpected error")
		t.FailNow()
	}
	doc.AppendChild(cmt)

	if cmt.GetNodeName() != "#comment" {
		t.Errorf("expecting '#comment', got '%v'", cmt.GetNodeName())
	}
	if cmt.GetNodeType() != CommentNode {
		t.Errorf("expecting CommentNode as type, but got %v", cmt.GetNodeType())
	}
	if cmt.GetNodeValue() != "comment string" {
		t.Error("node value should return 'comment string'")
	}
	if cmt.GetComment() != cmt.GetNodeValue() {
		t.Error("GetComment() should equal GetNodeValue()")
	}
	if cmt.GetParentNode() != doc {
		t.Error("got wrong parent node")
	}
	if cmt.GetAttributes() != nil {
		t.Error("comments can not have attributes and should therefore be nil")
	}
	if cmt.GetOwnerDocument() != doc {
		t.Error("invalid owner document")
	}
	if err := cmt.AppendChild(doc.CreateText("bleh")); err == nil {
		t.Error("appending a child to a comment should generate an error")
	}
	if cmt.HasChildNodes() != false {
		t.Error("should return false, was true")
	}
	if cmt.GetFirstChild() != nil {
		t.Error("first child should be nil at all times")
	}
	if cmt.GetNamespaceURI() != "" {
		t.Error("namespace URI should be empty")
	}
}
