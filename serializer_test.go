package dom

import (
	"strings"
	"testing"
)

func serializeToString(n Node) string {
	w := &strings.Builder{}

	ser := NewSerializer()
	ser.Configuration.NamespaceDeclarations = true
	ser.Configuration.PrettyPrint = true
	ser.Configuration.Namespaces = true
	ser.Serialize(n, w)

	return w.String()
}

func TestSerializationWithRootedNamespaceAndPrefix(t *testing.T) {
	doc := NewDocument()
	root, _ := doc.CreateElementNS("urn:doc", "ns1:rootElement")
	doc.AppendChild(root)

	childElement, _ := doc.CreateElement("ns1:childElement")
	childElement.SetTextContent("Text content")
	root.AppendChild(childElement)

	noNsWithPfx, _ := doc.CreateElement("nons:noNamespaceWithPrefix")
	noNsWithPfx.SetTextContent("This element has a prefix, but no namespace")
	root.AppendChild(noNsWithPfx)

	expected := `<?xml version="1.0" encoding="UTF-8"?>
<ns1:rootElement xmlns:ns1="urn:doc">
    <ns1:childElement>Text content</ns1:childElement>
    <nons:noNamespaceWithPrefix>This element has a prefix, but no namespace</nons:noNamespaceWithPrefix>
</ns1:rootElement>
`
	actual := serializeToString(doc)
	if expected != actual {
		t.Errorf("Expected:\n%s\nActual:\n%s", expected, actual)
	}
}

func TestSerializationTwiceSameNamespace(t *testing.T) {
	doc := NewDocument()
	root, _ := doc.CreateElementNS("urn:doc", "ns1:rootElement")
	doc.AppendChild(root)

	childElement, _ := doc.CreateElementNS("urn:doc", "ns1:childElement")
	childElement.SetTextContent("Text content")
	root.AppendChild(childElement)

	expected := `<?xml version="1.0" encoding="UTF-8"?>
<ns1:rootElement xmlns:ns1="urn:doc">
    <ns1:childElement>Text content</ns1:childElement>
</ns1:rootElement>
`
	actual := serializeToString(doc)
	if expected != actual {
		t.Errorf("Expected:\n%s\nActual:\n%s", expected, actual)
	}
}

func TestSerializationDefaultNamespacePrefix(t *testing.T) {
	doc := NewDocument()
	root, _ := doc.CreateElementNS("urn:doc", "rootElement")
	doc.AppendChild(root)

	childElement, _ := doc.CreateElement("childElement")
	childElement.SetTextContent("Text content")
	root.AppendChild(childElement)

	expected := `<?xml version="1.0" encoding="UTF-8"?>
<rootElement xmlns="urn:doc">
    <childElement>Text content</childElement>
</rootElement>
`

	actual := serializeToString(doc)
	if expected != actual {
		t.Errorf("Expected:\n%s\nActual:\n%s", expected, actual)
	}
}

func TestSerializationComments(t *testing.T) {
	doc := NewDocument()
	root, _ := doc.CreateElementNS("urn:doc", "rootElement")
	doc.AppendChild(root)

	childElement, _ := doc.CreateElement("childElement")
	childElement.SetTextContent("Text content")
	root.AppendChild(childElement)

	comment, _ := doc.CreateComment("some comment")
	root.AppendChild(comment)

	otherChild, _ := doc.CreateElement("moar")
	root.AppendChild(otherChild)

	expected := `<?xml version="1.0" encoding="UTF-8"?>
<rootElement xmlns="urn:doc">
    <childElement>Text content</childElement>
    <!-- some comment -->
    <moar/>
</rootElement>
`

	actual := serializeToString(doc)
	if expected != actual {
		t.Errorf("Expected:\n%s\nActual:\n%s", expected, actual)
	}
}

// FIXME: this needs to be fixed after I'm done with the regular serialization
// func TestSerializationAfterReading(t *testing.T) {
// 	doc := NewDocument()
// 	aanleverResponse, _ := doc.CreateElementNS("http://logius.nl/digipoort/koppelvlakservices/1.2/", "aanleverResponse")
// 	doc.AppendChild(aanleverResponse)

// 	kenmerk, _ := doc.CreateElement("kenmerk")
// 	kenmerk.SetTextContent("215222fb-13f4-4d9b-99a1-e61369e72acb")

// 	berichtsoort, _ := doc.CreateElement("berichtsoort")
// 	berichtsoort.SetTextContent("SBA_OB_2020")

// 	aanleverResponse.AppendChild(kenmerk)
// 	aanleverResponse.AppendChild(berichtsoort)

// 	t.Log(serializeToString(doc))

// 	reader := strings.NewReader(serializeToString(doc))
// 	parser := NewParser(reader)
// 	newdoc, _ := parser.Parse()

// 	// FIXME: this deser/ser has crappy results. MoveNamespaceToRoot also doesn't work properly!
// 	// MoveNamespacesToRoot(newdoc)

// 	PrintTree(doc, os.Stdout)

// 	t.Logf(serializeToString(newdoc))
// }

func TestSerializationWithPrefixedRootAndDoublyDefinedNamespaces(t *testing.T) {
	doc := NewDocument()

	rootElement, _ := doc.CreateElementNS("urn:dom:main", "ns0:RootElement")
	childElement, _ := doc.CreateElementNS("urn:dom:main", "ns0:ChildElement")
	childElement.SetAttribute("xmlns:ns0", "urn:dom:main")
	childElement.SetAttribute("xmlns:yep", "urn:this-should-be-serialized")
	childElement.SetAttribute("xmlns", "urn:this-should-be-serialized")
	childElement.SetAttribute("age", "55")

	otherNsElement, _ := doc.CreateElementNS("urn:some-other:element", "pfx:OtherNsElement")
	childOfOther, _ := doc.CreateElementNS("urn:some-other:element", "pfx:ChildOfOther")
	otherNsElement.AppendChild(childOfOther)

	doc.AppendChild(rootElement)
	rootElement.AppendChild(childElement)
	rootElement.AppendChild(otherNsElement)

	expected := `<?xml version="1.0" encoding="UTF-8"?>
<ns0:RootElement xmlns:ns0="urn:dom:main">
    <ns0:ChildElement xmlns:yep="urn:this-should-be-serialized" xmlns="urn:this-should-be-serialized" age="55"/>
    <pfx:OtherNsElement xmlns:pfx="urn:some-other:element">
        <pfx:ChildOfOther/>
    </pfx:OtherNsElement>
</ns0:RootElement>
`
	actual := serializeToString(doc)

	if actual != expected {
		t.Errorf("Expected:\n%s\nActual:\n%s", expected, actual)
	}
}

func TestSerializationWithNonPrefixedRootAndDoublyDefinedNamespaces(t *testing.T) {
	doc := NewDocument()

	rootElement, _ := doc.CreateElementNS("urn:dom:main", "RootElement")
	childElement, _ := doc.CreateElementNS("urn:dom:main", "ChildElement")
	childElement.SetAttribute("xmlns", "urn:dom:main") // this should be omitted during the serialization
	childElement.SetAttribute("xmlns:yep", "urn:this-should-be-serialized")
	childElement.SetAttribute("age", "55")

	grandChild, _ := doc.CreateElementNS("urn:dom:main", "GrandChildElement")

	childElement.AppendChild(grandChild)

	otherNsElement, _ := doc.CreateElementNS("urn:some-other:element", "pfx:OtherNsElement")
	childOfOther, _ := doc.CreateElementNS("urn:some-other:element", "pfx:OtherChild")
	otherNsElement.AppendChild(childOfOther)

	doc.AppendChild(rootElement)
	rootElement.AppendChild(childElement)
	rootElement.AppendChild(otherNsElement)

	expected := `<?xml version="1.0" encoding="UTF-8"?>
<RootElement xmlns="urn:dom:main">
    <ChildElement xmlns:yep="urn:this-should-be-serialized" age="55">
        <GrandChildElement/>
    </ChildElement>
    <pfx:OtherNsElement xmlns:pfx="urn:some-other:element">
        <pfx:OtherChild/>
    </pfx:OtherNsElement>
</RootElement>
`
	actual := serializeToString(doc)

	if actual != expected {
		t.Errorf("Expected:\n%s\nActual:\n%s", expected, actual)
	}
}
