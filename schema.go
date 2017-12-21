package atomfeed

import "encoding/xml"

// Feed is an atom:feed element and is the document (i.e., top-level) element of
// an Atom Feed Document, acting as a container for metadata and data associated with the feed.
//  https://tools.ietf.org/html/rfc4287#section-4.1.1
type Feed struct {
	XMLName     xml.Name       `xml:"feed"`
	Namespace   string         `xml:"xmlns,attr"` // xmlns="http://www.w3.org/2005/Atom"
	ID          ID             `xml:"id"`
	Generator   *Generator     `xml:"generator"`
	Links       []Link         `xml:"link"`
	Updated     *Date          `xml:"updated"`
	Title       *TextConstruct `xml:"title"`
	Subtitle    *TextConstruct `xml:"subtitle"`
	Icon        *Icon          `xml:"icon"`
	Logo        *Logo          `xml:"logo"`
	Categories  []Category     `xml:"category"`
	Author      *Person        `xml:"author"`
	Contributor []Person       `xml:"contributor"`
	Copyright   *TextConstruct `xml:"rights"` // https://tools.ietf.org/html/rfc4287#section-4.2.10
	Entries     []Entry        `xml:"entry"`
	*CommonAttributes
}

// Entry is an atom:entry element and represents an individual entry, acting as a
// container for metadata and data associated with the entry.
//  https://tools.ietf.org/html/rfc4287#section-4.1.2
type Entry struct {
	ID          ID             `xml:"id"`
	Title       *TextConstruct `xml:"title"`
	Links       []Link         `xml:"link"`
	Published   *Date          `xml:"published"`
	Updated     *Date          `xml:"updated"`
	Author      *Person        `xml:"author"`
	Categories  []Category     `xml:"category"`
	Copyright   *TextConstruct `xml:"rights"`
	Contributor []Person       `xml:"contributor"`
	Source      *Source        `xml:"source"`
	Summary     *Content       `xml:"summary"`
	Content     *Content       `xml:"content"`
	*CommonAttributes
}

// Source is an atom:source element.
//  https://tools.ietf.org/html/rfc4287#section-4.2.11
type Source struct {
	ID          *ID            `xml:"id"`
	Generator   *Generator     `xml:"generator"`
	Links       []Link         `xml:"link"`
	Updated     *Date          `xml:"updated"`
	Title       *TextConstruct `xml:"title"`
	Subtitle    *TextConstruct `xml:"subtitle"`
	Icon        *Icon          `xml:"icon"`
	Logo        *Logo          `xml:"logo"`
	Categories  []Category     `xml:"category"`
	Author      *Person        `xml:"author"`
	Contributor []Person       `xml:"contributor"`
	Copyright   *TextConstruct `xml:"rights"`
	*CommonAttributes
}

// Content is an atom:content element which either contains or links to the content of the entry.
//  https://tools.ietf.org/html/rfc4287#section-4.1.3
type Content struct {
	// Type MAY be one of "text", "html", or "xhtml".
	// https://tools.ietf.org/html/rfc4287#section-4.1.3.1
	Type string `xml:"type,attr,omitempty"`
	// Source is an optional attribute, whose value MUST be an IRI reference [RFC3987].
	// If the "src" attribute is present, atom:content MUST be empty.
	// https://tools.ietf.org/html/rfc4287#section-4.1.3.2
	Source string `xml:"src,attr,omitempty"`
	Value  string `xml:",chardata"`
	*CommonAttributes
	base64Encoded bool
}

// TextConstruct contains human-readable text, usually in small quantities.
//  https://tools.ietf.org/html/rfc4287#section-3.1
type TextConstruct struct {
	// Type MAY be one of "text", "html", or "xhtml".
	// https://tools.ietf.org/html/rfc4287#section-3.1.1
	Type  string `xml:"type,attr,omitempty"`
	Value string `xml:",chardata"`
	*CommonAttributes
}

// Person is an atom:person element that describes a person, corporation, or similar entity.
// Use the person struct to create author, coauthor or contributor entities.
//  https://tools.ietf.org/html/rfc4287#section-3.2
type Person struct {
	Name  string `xml:"name"`
	Email string `xml:"email,omitempty"`
	URI   string `xml:"uri,omitempty"`
	*CommonAttributes
}

// Date is an atom:date element whose content MUST conform to the "date-time" format defined in [RFC3339].
//  https://tools.ietf.org/html/rfc4287#section-3.3
type Date struct {
	// Value must be a date conforming to RFC3339.
	// Try: `time.Now().Format(time.RFC3339)`
	Value string `xml:",chardata"`
	*CommonAttributes
}

// Generator is an atom:generator element.
// The "atom:generator" element's content identifies the agent used to generate a feed.
//  https://tools.ietf.org/html/rfc4287#section-4.2.4
type Generator struct {
	URI     string `xml:"uri,attr,omitempty"`
	Version string `xml:"version,attr,omitempty"`
	Value   string `xml:",chardata"`
	*CommonAttributes
}

// Category is an atom:category element.
//  https://tools.ietf.org/html/rfc4287#section-4.2.2
type Category struct {
	// Term is a mandatory string that identifies the category
	// to which the entry or feed belongs.
	// https://tools.ietf.org/html/rfc4287#section-4.2.2.1
	Term string `xml:"term,attr"`
	// Scheme is an optional IRI that identifies a categorization scheme.
	// https://tools.ietf.org/html/rfc4287#section-4.2.2.2
	Scheme string `xml:"scheme,attr,omitempty"`
	// Label provides an optional human-readable label for display in end-user applications.
	// https://tools.ietf.org/html/rfc4287#section-4.2.2.3
	Label string `xml:"label,attr,omitempty"`
	*CommonAttributes
}

// Icon is an optional atom:icon element, which is an IRI reference [RFC3987] that
// identifies an image that provides iconic visual identification for a feed.
//
// The image SHOULD have an aspect ratio of 1 (horizontal) to 1 (vertical)
// and SHOULD be suitable for presentation at a small size.
//  https://tools.ietf.org/html/rfc4287#section-4.2.5
type Icon struct {
	Value string `xml:",chardata"`
	*CommonAttributes
}

// ID is an atom:id element and conveys a permanent, universally unique identifier for an entry or feed.
//
// The value of ID must be a valid IRI.
// A permalink SHOULDN'T be used as an ID.
//  https://github.com/denisbrodbeck/atomfeed/blob/master/README.md#id
//  https://tools.ietf.org/html/rfc4287#section-4.2.6
type ID struct {
	Value string `xml:",chardata"`
	*CommonAttributes
}

// Link is an atom:link element that defines a reference from an entry or feed to a Web resource.
//  https://tools.ietf.org/html/rfc4287#section-4.2.7
type Link struct {
	// Href contains the link's mandatory IRI.
	// https://tools.ietf.org/html/rfc4287#section-4.2.7.1
	Href string `xml:"href,attr"`
	// Rel is an optional attribute that indicates the link relation type.
	// https://tools.ietf.org/html/rfc4287#section-4.2.7.2
	Rel string `xml:"rel,attr,omitempty"`
	// Type is an optional advisory media type.
	// It is a hint about the type of the representation that is
	// expected to be returned when the value of the href attribute is
	// dereferenced.  Note that the type attribute does not override the
	// actual media type returned with the representation.
	// https://tools.ietf.org/html/rfc4287#section-4.2.7.3
	Type string `xml:"type,attr,omitempty"`
	// HrefLang describes the optional language of the resource pointed to by the href attribute.
	// When used together with the rel="alternate", it implies a translated version of the entry.
	// https://tools.ietf.org/html/rfc4287#section-4.2.7.4
	HrefLang string `xml:"hreflang,attr,omitempty"`
	// Title conveys optional human-readable information about the link.
	// https://tools.ietf.org/html/rfc4287#section-4.2.7.5
	Title string `xml:"title,attr,omitempty"`
	// Length indicates an optional advisory length of the linked content in octets.
	// https://tools.ietf.org/html/rfc4287#section-4.2.7.6
	Length string `xml:"length,attr,omitempty"`
	*CommonAttributes
}

// Logo is an atom:logo element, which is n IRI reference [RFC3987] that
// identifies an image that provides visual identification for a feed.
//
// The image SHOULD have an aspect ratio of 2 (horizontal) to 1 (vertical).
//  https://tools.ietf.org/html/rfc4287#section-4.2.8
type Logo struct {
	Value string `xml:",chardata"`
	*CommonAttributes
}

// CommonAttributes is an atom:commonattributes element.
//  https://tools.ietf.org/html/rfc4287#section-2
type CommonAttributes struct {
	Base string `xml:"xml:base,attr,omitempty"`
	Lang string `xml:"xml:lang,attr,omitempty"`
}
