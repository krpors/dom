package dom

import (
	"testing"
)

func TestAttrNodeName(t *testing.T) {
	attr := newAttr()
	attr.setName("cruft")
	if attr.GetName() != "cruft" {
		t.Errorf("incorrect node name")
	}

	if attr.GetNodeName() != "cruft" {
		t.Errorf("incorrect node name")
	}
}

func TestAttrNodeType(t *testing.T) {
	attr := newAttr()
	if attr.GetNodeType() != AttributeNode {
		t.Errorf("incorrect node type for attribute")
	}
}

func TestAttrNodeValue(t *testing.T) {
	attr := newAttr()
	attr.SetValue("valval")
	if attr.GetNodeValue() != "valval" {
		t.Errorf("incorrect node value: '%v'", attr.GetNodeValue())
	}

	if attr.GetValue() != "valval" {
		t.Errorf("incorrect node value: '%v'", attr.GetValue())
	}
}

func TestAttrParentNode(t *testing.T) {
	elem := newElement()
	elem.SetTagName("root")

	attr := newAttr()
	attr.setName("attrname")
	attr.SetValue("attrvalue")

	attr.setParentNode(elem)
}

func TestAttrAppendChild(t *testing.T) {
	attr := newAttr()
	comment := newComment()
	comment.SetComment("hi thar")
	if err := attr.AppendChild(comment); err == nil {
		t.Errorf("expected error during appending of node to attr")
	}
}
