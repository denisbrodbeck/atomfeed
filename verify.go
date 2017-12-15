package atomfeed

import (
	"fmt"
	"net/mail"
	"net/url"
	"strings"
)

// Verify checks an atom:feed element for most common errors.
//
// Common checks are the existence of atom:id, atom:author,
// atom:title, atom:updated. Any entries will be checked, too.
func (f *Feed) Verify() error {
	if err := checkID(f.ID); err != nil {
		return err
	}
	if err := checkAuthorsExist(f); err != nil {
		return err
	}
	if err := checkPerson(f.Author); err != nil {
		return err
	}
	if f.Logo != nil {
		if err := checkURI(f.Logo.Value); err != nil {
			return err
		}
	}
	if f.Icon != nil {
		if err := checkURI(f.Icon.Value); err != nil {
			return err
		}
	}
	if f.Title == nil || f.Title.Value == "" {
		return fmt.Errorf("atom:feed elements needs a title")
	}
	if f.Updated == nil {
		return fmt.Errorf("atom:feed elements needs a atom:updated")
	}
	for _, entry := range f.Entries {
		if err := entry.Verify(); err != nil {
			return err
		}
	}
	return nil
}

// Verify checks an atom:entry element for most common errors.
//
// Common checks are the existence of atom:id, atom:author,
// atom:title, atom:updated, atom:content.
func (e *Entry) Verify() error {
	if err := checkID(e.ID); err != nil {
		return err
	}
	if err := checkContent(e.Content); err != nil {
		return err
	}
	if err := checkPerson(e.Author); err != nil {
		return err
	}
	if e.Title == nil || e.Title.Value == "" {
		return fmt.Errorf("atom:entry elements need an atom:title element")
	}
	if e.Updated == nil {
		return fmt.Errorf("atom:entry elements need an atom:updated element")
	}
	if e.Content.Source != "" {
		if e.Summary.Value == "" {
			return fmt.Errorf("atom:entry elements need an atom:summary element because atom:content element has a src attribute set")
		}
	} else if e.Content.base64Encoded {
		if e.Summary.Value == "" {
			return fmt.Errorf("atom:entry elements need an atom:summary element because atom:content is base64 encoded")
		}
	}
	return nil
}

func checkAuthorsExist(f *Feed) error {
	hasFeedAuthor := f.Author != nil && f.Author.Name != ""
	if hasFeedAuthor == false {
		if len(f.Entries) == 0 {
			return fmt.Errorf("missing author field: an atom:feed must have an atom:author unless all of its atom:entry children have an atom:author")
		}
		allEntriesHaveAuthor := true
		for _, entry := range f.Entries {
			if entry.Author == nil || entry.Author.Name == "" {
				allEntriesHaveAuthor = false
				break
			}
		}
		if allEntriesHaveAuthor == false {
			return fmt.Errorf("missing author field: an atom:feed must have an atom:author unless all of its atom:entry children have an atom:author")
		}
	}
	return nil
}

func checkPerson(p *Person) error {
	if p == nil {
		return nil
	}
	if p.Name == "" {
		return fmt.Errorf("author name cannot be empty")
	}
	if err := checkEmail(p.Email); err != nil {
		return err
	}
	return checkURI(p.URI)
}

func checkEmail(email string) error {
	if email == "" {
		return nil
	}
	if _, err := mail.ParseAddress(email); err != nil {
		return fmt.Errorf("%q is not a valid email address: %v", email, err)
	}
	return nil
}

func checkURI(uri string) error {
	if uri == "" {
		return nil
	}
	if _, err := url.Parse(uri); err != nil {
		return fmt.Errorf("%q is not a valid URI: %v", uri, err)
	}
	return nil
}

func checkID(id *ID) error {
	if id == nil || id.Value == "" {
		return fmt.Errorf("ID cannot be empty")
	}
	return checkURI(id.Value)
}

func checkContent(c *Content) error {
	if c.Source != "" { // https://tools.ietf.org/html/rfc4287#section-4.1.3.2
		if err := checkURI(c.Source); err != nil { // MUST be IRI
			return err
		}
		if c.Value != "" { // MUST be empty
			return fmt.Errorf("invalid content: src attribute is present, therefore content must be empty")
		}
		if c.Type != "" { // SHOULD be provided
			if strings.Contains(c.Type, "/") == false { // MUST be mime
				return fmt.Errorf("invalid mime type: %v", c.Type)
			}
		}
	}
	switch c.Type {
	case "", "text", "html", "xhtml":
		return nil
	default:
		// Whatever a media type is, it contains at least one slash
		if strings.Contains(c.Type, "/") == false {
			return fmt.Errorf("invalid mime type: %v", c.Type)
		}
	}
	return nil
}
