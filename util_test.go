package dom

import (
	"bytes"
	"strings"
	"testing"
)

// TestUtilEscape tests the convenience method to escape XML character data.
func TestUtilEscape(t *testing.T) {
	data := "<>&&amp;"
	str := escape(data)
	expected := "&lt;&gt;&amp;&amp;amp;"
	if str != expected {
		t.Logf("expected: %v", expected)
		t.Logf("acuta:    %v", str)
		t.Errorf("escaped sequence does not match")
	}
}

func TestToXML(t *testing.T) {
	doc := NewDocument()

	procinst, _ := doc.CreateProcessingInstruction("violin", "tdavis")
	root, _ := doc.CreateElement("root")
	whitespaceText := doc.CreateText("  ")
	comment, _ := doc.CreateComment("rox your sox")
	first, _ := doc.CreateElement("first")
	text := doc.CreateText("< & > are entities!")
	second, _ := doc.CreateElement("second")
	nochildren, _ := doc.CreateElement("nochildren")
	attr, _ := doc.CreateAttribute("name")
	attr.SetValue("value")

	doc.AppendChild(procinst)
	doc.AppendChild(root)
	root.AppendChild(whitespaceText)
	root.AppendChild(first)
	root.AppendChild(comment)
	root.AppendChild(second)

	first.AppendChild(text)
	second.AppendChild(nochildren)
	nochildren.SetAttributeNode(attr)

	var buf bytes.Buffer
	ToXML(doc, false, &buf)

	bla := []string{
		`<?xml version="1.0" encoding="UTF-8"?>`,
		`<?violin tdavis?>`,
		`<root>  <first>&lt; &amp; &gt; are entities!</first><!-- rox your sox --><second><nochildren name="value"/></second></root>`,
	}

	expected := strings.Join(bla, "")

	if buf.String() != expected {
		t.Logf("actual:   %v", buf.String())
		t.Logf("expected: %v", expected)
		t.Fail()
	}
}

func TestUtilToXML2(t *testing.T) {
	doc := NewDocument()

	contacts, _ := doc.CreateElementNS("urn:contacts:contacts", "contacts")
	contact, _ := doc.CreateElementNS("urn:contacts:contact", "contact")
	person, _ := doc.CreateElementNS("urn:contacts:person", "person")
	name, _ := doc.CreateElement("name")
	text := doc.CreateText("some text.")

	doc.AppendChild(contacts)
	contacts.AppendChild(contact)
	contact.AppendChild(person)
	person.AppendChild(name)
	name.AppendChild(text)

	//	ToXML(doc, false, os.Stdout)
}
