package dom

import (
	"fmt"
	"os"
	"strings"
	"testing"
)

var exampleDoc = `<?xml version="1.0" encoding="UTF-8"?>
<directory>
	<people>
		<person name="Foo" lastname="Quux">
			<date>2016-01-01</date>
		</person>
		<person name="Alice" lastname="Bob">
			<date>2016-09-03</date>
		</person>
		<!-- this is a comment -->
		<!-- empty element follows -->
		<person/>
	</people>
	<ns:cruft xmlns:ns="http://exmple.org/xmlns/uri">
		<ns:other>Character data.</ns:other>
		<ns:balls ns:derp="woot">More chardata</ns:balls>
	</ns:cruft>
</directory>
`

func TestWut(t *testing.T) {
	reader := strings.NewReader(exampleDoc)
	builder := NewBuilder(reader)
	doc, _ := builder.CreateDocument()

	fmt.Println("===========")
	builder.PrintTree(os.Stdout)
	fmt.Println("===========")
	fmt.Println(ToXML(doc))

	t.Fail()

	// asd
}
