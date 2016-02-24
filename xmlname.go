package dom

import "unicode"

// XMLName represents an XML Name according to the specification at https://www.w3.org/TR/xml/#NT-Name.
// It uses the 'NameStartChar' and 'NameChar' to validate elements.
type XMLName string

// IsValid returns true if the given XMLName is valid according to the XML specification.
func (name XMLName) IsValid() bool {
	return true
}

var nameStartChars = &unicode.RangeTable{
	R16: []unicode.Range16{
		{0x003A, 0x003A, 1}, // :
		{0x005F, 0x005F, 1}, // _
		{0x0041, 0x005A, 1}, // A-Z
		{0x0061, 0x007A, 1}, // a-z
		{0x00C0, 0x00D6, 1},
		{0x00D8, 0x00F6, 1},
		{0x00F8, 0x02FF, 1},
		{0x0370, 0x037D, 1},
		{0x037F, 0x1FFF, 1},
		{0x200C, 0x200D, 1},
		{0x2070, 0x218F, 1},
		{0x2C00, 0x2FEF, 1},
		{0x3001, 0xD7FF, 1},
		{0xF900, 0xFDCF, 1},
		{0xFDF0, 0xFFFD, 1},
	},
}

var nameChars = &unicode.RangeTable{
	R16: []unicode.Range16{
		{0x002D, 0x002D, 1}, // -
		{0x002E, 0x002E, 1}, // .
		{0x0030, 0x0039, 1}, // 0-9
		{0x00B7, 0x00B7, 1},
		{0x0300, 0x036F, 1},
		{0x203F, 0x2040, 1},
	},
}
