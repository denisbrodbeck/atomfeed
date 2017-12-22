package atomfeed

import (
	"fmt"
	"net/mail"
	"net/url"
	"strings"
	"time"
)

// VerificationError describes problems encountered during feed verification.
type VerificationError struct {
	Errors []error
}

func (e *VerificationError) Error() string {
	errors := []string{}
	for _, e := range e.Errors {
		errors = append(errors, e.Error())
	}
	return strings.Join(errors, "\n")
}

// Verify checks an atom:feed element for most common errors.
//
// Common checks are the existence of atom:id, atom:author,
// atom:title, atom:updated. Any entries will be checked, too.
func (f *Feed) Verify() *VerificationError {
	errors := []error{}
	if err := checkID(f.ID); err != nil {
		errors = append(errors, fmt.Errorf("feed: %v", err))
	}
	if err := checkAuthorsExist(f); err != nil {
		errors = append(errors, fmt.Errorf("feed: %v", err))
	}
	if err := checkPerson(f.Author); err != nil {
		errors = append(errors, fmt.Errorf("feed: author: %v", err))
	}
	if f.Logo != nil {
		if err := checkURI(f.Logo.Value); err != nil {
			errors = append(errors, fmt.Errorf("feed: logo: %v", err))
		}
	}
	if f.Icon != nil {
		if err := checkURI(f.Icon.Value); err != nil {
			errors = append(errors, fmt.Errorf("feed: icon: %v", err))
		}
	}
	if f.Title == nil || f.Title.Value == "" {
		errors = append(errors, fmt.Errorf("feed: missing title"))
	}
	if f.Updated == nil {
		errors = append(errors, fmt.Errorf("feed: missing updated date"))
	} else {
		if err := checkDate(f.Updated.Value); err != nil {
			errors = append(errors, fmt.Errorf("feed: updated: %v", err))
		}
	}
	for _, entry := range f.Entries {
		if err := entry.Verify(); err != nil {
			errors = append(errors, fmt.Errorf("errors for entry [%s]\n%v", entry.String(), err))
		}
	}
	if len(errors) > 0 {
		return &VerificationError{Errors: errors}
	}
	return nil
}

// Verify checks an atom:entry element for most common errors.
//
// Common checks are the existence of atom:id, atom:author,
// atom:title, atom:updated, atom:content.
func (e *Entry) Verify() *VerificationError {
	errors := []error{}
	if err := checkID(e.ID); err != nil {
		errors = append(errors, fmt.Errorf("entry: %v", err))
	}
	if err := checkContent(e.Content); err != nil {
		errors = append(errors, fmt.Errorf("entry: %v", err))
	}
	if err := checkPerson(e.Author); err != nil {
		errors = append(errors, fmt.Errorf("entry: author: %v", err))
	}
	if e.Title == nil || e.Title.Value == "" {
		errors = append(errors, fmt.Errorf("entry: missing title"))
	}
	if e.Updated == nil {
		errors = append(errors, fmt.Errorf("entry: missing updated date"))
	} else {
		if err := checkDate(e.Updated.Value); err != nil {
			errors = append(errors, fmt.Errorf("entry: updated: %v", err))
		}
	}
	if e.Published != nil {
		if err := checkDate(e.Published.Value); err != nil {
			errors = append(errors, fmt.Errorf("entry: published: %v", err))
		}
	}
	if e.Content != nil {
		if e.Content.Source != "" {
			if e.Summary == nil || e.Summary.Value == "" {
				errors = append(errors, fmt.Errorf("entry: need a summary because content has src attribute set"))
			}
		} else if e.Content.base64Encoded {
			if e.Summary == nil || e.Summary.Value == "" {
				errors = append(errors, fmt.Errorf("entry: need a summary because content is base64 encoded"))
			}
		}
	}
	if len(errors) > 0 {
		return &VerificationError{Errors: errors}
	}
	return nil
}

func checkAuthorsExist(f *Feed) error {
	hasFeedAuthor := f.Author != nil && f.Author.Name != ""
	if hasFeedAuthor == false {
		if len(f.Entries) == 0 {
			return fmt.Errorf("missing author field: an atom feed must have an author unless all of its entry children have an author")
		}
		allEntriesHaveAuthor := true
		for _, entry := range f.Entries {
			if entry.Author == nil || entry.Author.Name == "" {
				allEntriesHaveAuthor = false
				break
			}
		}
		if allEntriesHaveAuthor == false {
			return fmt.Errorf("missing author field: an atom feed must have an author unless all of its entry children have an author")
		}
	}
	return nil
}

func checkPerson(p *Person) error {
	if p == nil {
		return nil
	}
	if p.Name == "" {
		return fmt.Errorf("name cannot be empty")
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

func checkID(id ID) error {
	if id.Value == "" {
		return fmt.Errorf("ID cannot be empty")
	}
	return checkURI(id.Value)
}

func checkDate(date string) error {
	t, err := time.Parse(time.RFC3339, date)
	if err != nil {
		return fmt.Errorf("invalid date: %v", err)
	}
	if t.IsZero() {
		return fmt.Errorf("invalid date %q: date is zero", date)
	}
	return nil
}

func checkContent(c *Content) error {
	if c == nil {
		return nil
	}
	if c.Source != "" { // https://tools.ietf.org/html/rfc4287#section-4.1.3.2
		if err := checkURI(c.Source); err != nil { // MUST be IRI
			return err
		}
		if c.Value != "" || c.ValueXML != "" { // MUST be empty
			return fmt.Errorf("invalid content: src attribute is present, therefore content must be empty")
		}
		if c.Type != "" { // SHOULD be provided
			if strings.Contains(c.Type, "/") == false { // MUST be mime
				return fmt.Errorf("invalid mime type: %v", c.Type)
			}
		}
	}
	switch c.Type {
	case "", "text", "html":
		if c.ValueXML != "" {
			return fmt.Errorf("field %q must be empty when using type %q — use field %q instead", "ValueXML", c.Type, "Value")
		}
	case "xhtml":
		if c.Value != "" {
			return fmt.Errorf("field %q must be empty when using type %q — use field %q instead", "Value", "xhtml", "ValueXML")
		}
	default:
		// Whatever a media type is, it contains at least one slash
		if strings.Contains(c.Type, "/") == false {
			return fmt.Errorf("invalid mime type: %v", c.Type)
		}
	}
	return nil
}
