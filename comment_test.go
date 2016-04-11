package dom

import (
	"testing"
)

// Tests the getters of a Comment node.
func TestCommentGetters(t *testing.T) {
	doc := NewDocument()
	cmt1, err := doc.CreateComment("comment string")
	if err != nil {
		t.Error("unexpected error")
		t.FailNow()
	}
	doc.AppendChild(cmt1)

	if cmt1.GetNodeName() != "#comment" {
		t.Errorf("expecting '#comment', got '%v'", cmt1.GetNodeName())
	}
	if cmt1.GetNodeType() != CommentNode {
		t.Errorf("expecting CommentNode as type, but got %v", cmt1.GetNodeType())
	}
	if cmt1.GetNodeValue() != "comment string" {
		t.Error("node value should return 'comment string'")
	}
	if cmt1.GetComment() != cmt1.GetNodeValue() {
		t.Error("GetComment() should equal GetNodeValue()")
	}
	if cmt1.GetParentNode() != doc {
		t.Error("got wrong parent node")
	}
	if cmt1.GetAttributes() != nil {
		t.Error("comments can not have attributes and should therefore be nil")
	}
	if cmt1.GetOwnerDocument() != doc {
		t.Error("invalid owner document")
	}
	if err := cmt1.AppendChild(doc.CreateText("bleh")); err == nil {
		t.Error("appending a child to a comment should generate an error")
	}
	if _, err := cmt1.RemoveChild(nil); err == nil {
		t.Error("removing a child should generate an error")
	}
	if _, err := cmt1.InsertBefore(nil, nil); err == nil {
		t.Error("inserting a child should generate an error")
	}
	if _, err := cmt1.ReplaceChild(nil, nil); err == nil {
		t.Error("replacing a child should generate an error")
	}
	if cmt1.HasChildNodes() != false {
		t.Error("should return false, was true")
	}
	if cmt1.GetFirstChild() != nil {
		t.Error("first child should be nil at all times")
	}
	if cmt1.GetLastChild() != nil {
		t.Error("last child should be nil at all times")
	}
	if cmt1.GetNamespaceURI() != "" {
		t.Error("namespace URI should be empty")
	}
	if cmt1.GetLocalName() != "" {
		t.Error("locale name should be an empty string")
	}
	if cmt1.GetNamespacePrefix() != "" {
		t.Error("comments should not have namespace prefixes")
	}
	if cmt1.HasAttributes() {
		t.Error("comments cannot have attributes")
	}
}

// Tests the appending of nodes, and getting previous and next siblings.
func TestCommentSiblings(t *testing.T) {
	doc := NewDocument()
	cmt1, _ := doc.CreateComment("comment 1")
	cmt2, _ := doc.CreateComment("comment 2")
	cmt3, _ := doc.CreateComment("comment 3")

	doc.AppendChild(cmt1)
	doc.AppendChild(cmt2)
	doc.AppendChild(cmt3)

	sibling := cmt1.GetNextSibling()
	if sibling != cmt2 {
		t.Errorf("next sibling of 'comment' should be 'comment 2', but was '%v'", sibling)
	}
	sibling = sibling.GetNextSibling()
	if sibling != cmt3 {
		t.Errorf("next sibling of 'comment 2' should be 'comment 3', but was '%v'", sibling)
	}
	sibling = sibling.GetPreviousSibling()
	if sibling != cmt2 {
		t.Errorf("previous sibling of 'comment 3' should be 'comment 2', but was '%v'", sibling)
	}
}
