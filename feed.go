package atomfeed // import "github.com/denisbrodbeck/atomfeed"

import (
	"encoding/base64"
	"encoding/xml"
	"fmt"
	"io"
	"strings"
	"time"
)

// Encode writes the XML encoding of Feed to the stream.
func (f *Feed) Encode(w io.Writer) error {
	enc := xml.NewEncoder(w)
	enc.Indent("", "  ")
	w.Write([]byte(xml.Header))
	return enc.Encode(f)
}

// NewFeed creates a basic atom:feed element suitable for e.g. a blog.
func NewFeed(id ID, author *Person, title, subtitle, baseURL, feedURL string, updated time.Time, entries []Entry) Feed {
	generator := &Generator{
		URI:     "https://github.com/denisbrodbeck/atomfeed",
		Version: "1.0",
		Value:   "atomfeed package",
	}
	return Feed{
		Namespace: "http://www.w3.org/2005/Atom",
		ID:        id,
		Title:     &TextConstruct{Value: title},
		Subtitle:  &TextConstruct{Value: subtitle},
		Author:    author,
		Links: []Link{
			{
				Rel:  "alternate",
				Type: "text/html",
				Href: baseURL, // https://example.com/
			},
			{
				Rel:  "self",
				Type: "application/atom+xml",
				Href: feedURL, // https://example.com/feed.atom
			},
		},
		Updated:   NewDate(updated),
		Entries:   entries,
		Generator: generator,
	}
}

// NewID creates an atom:id element.
// The given id parameter is taken as a raw value for ID.
func NewID(id string) ID {
	return ID{Value: id}
}

// NewFeedID creates a stable ID for an atom:feed element.
// The resulting ID follows the 'tag' URI scheme as defined in RFC 4151.
// More specifically the function creates valid atom IDs by feed creation time and a custom specifier.
//
// Further info:
//  https://github.com/denisbrodbeck/atomfeed/blob/master/README.md#id
//  http://web.archive.org/web/20110514113830/http://diveintomark.org/archives/2004/05/28/howto-atom-id
//  https://tools.ietf.org/html/rfc4151
func NewFeedID(authorityName string, creationTime time.Time, specific string) ID {
	tag := fmt.Sprintf("tag:%s,%s:%s", authorityName, creationTime.Format("2006-01-02"), specific)
	return ID{Value: tag}
}

// NewEntryID creates a stable ID for an atom:entry element.
// The resulting ID follows the 'tag' URI scheme as defined in RFC 4151.
// More specifically the function creates valid atom IDs by article creation time.
//
// Further info:
//  https://github.com/denisbrodbeck/atomfeed/blob/master/README.md#id
//  http://web.archive.org/web/20110514113830/http://diveintomark.org/archives/2004/05/28/howto-atom-id
//  https://tools.ietf.org/html/rfc4151
func NewEntryID(feedID ID, entryCreationTime time.Time) ID {
	tag := fmt.Sprintf("%s.post-%s", feedID.Value, entryCreationTime.Format("20060102150405"))
	return ID{Value: tag}
}

// NewContent creates the correct atom:content element depending on type attribute.
//
// https://tools.ietf.org/html/rfc4287#section-4.1.3.3
func NewContent(contentType, source string, value []byte) *Content {
	if source == "" && (value == nil || len(value) == 0) {
		return nil
	}
	switch {
	case contentType == "xhtml",
		contentType == "text/xml", // https://tools.ietf.org/html/rfc3023#section-3
		contentType == "application/xml",
		contentType == "text/xml-external-parsed-entity",
		contentType == "application/xml-external-parsed-entity",
		contentType == "application/xml-dtd",
		strings.HasSuffix(strings.ToLower(contentType), "+xml"),
		strings.HasSuffix(strings.ToLower(contentType), "/xml"):
		return &Content{Type: contentType, Source: source, ValueXML: string(value)}
	case contentType == "",
		contentType == "text",
		contentType == "html",
		strings.HasPrefix(strings.ToLower(contentType), "text/"):
		return &Content{Type: contentType, Source: source, Value: string(value)}
	}
	// all other types MUST be base64 encoded
	return &Content{Type: contentType, Source: source, Value: base64.StdEncoding.EncodeToString(value), base64Encoded: true}
}

// NewDate returns an atom:date element with valid RFC3339 time data.
func NewDate(t time.Time) *Date {
	if t.IsZero() {
		return nil
	}
	return &Date{Value: t.Format(time.RFC3339)}
}

// NewPerson returns an atom:person element.
func NewPerson(name, email, uri string) *Person {
	return &Person{Name: name, Email: email, URI: uri}
}

// NewCategory returns an atom:category element.
func NewCategory(category string) *Category {
	return &Category{Term: category}
}

// NewEntry creates a basic atom:entry suitable for e.g. a blog.
func NewEntry(id ID, title, permalink string, author *Person, updated, published time.Time, categories []string, summary, content []byte) Entry {
	return Entry{
		ID:    id,
		Title: &TextConstruct{Value: title},
		Links: []Link{
			{
				Rel:  "alternate",
				Type: "text/html",
				Href: permalink,
			},
		},
		Published:  NewDate(published),
		Updated:    NewDate(updated),
		Author:     author,
		Categories: termsToCategories(categories),
		Summary:    NewContent("html", "", summary),
		Content:    NewContent("html", "", content),
	}
}

func termsToCategories(categories []string) []Category {
	cat := []Category{}
	for _, c := range categories {
		cat = append(cat, *NewCategory(c))
	}
	return cat
}

func (e *Entry) String() string {
	title := ""
	if e.Title != nil {
		title = e.Title.Value
	}
	updated := ""
	if e.Updated != nil {
		updated = e.Updated.Value
	}
	published := ""
	if e.Published != nil {
		published = e.Published.Value
	}
	author := ""
	if e.Author != nil {
		author = e.Author.Name
	}
	categories := []string{}
	for _, c := range e.Categories {
		categories = append(categories, c.Term)
	}
	return fmt.Sprintf(
		"ID: %q, Title: %q, Updated: %q, Published: %q, Author: %q, Categories: %q",
		e.ID.Value,
		title,
		updated,
		published,
		author,
		strings.Join(categories, ","),
	)
}
