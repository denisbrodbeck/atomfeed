package atomfeed

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"log"
	"reflect"
	"testing"
	"time"
)

func TestNewFeed(t *testing.T) {
	authority := "example.com"
	now := time.Date(2012, time.December, 21, 8, 30, 15, 0, time.UTC)
	feedID := NewID(fmt.Sprintf("tag:%s,%s:blog", authority, now.Format("2006-01-02")))
	title := "example.com blog"
	subtitle := "Get the very latest news from the net."
	author := NewPerson("Go Pher", "", "https://blog.golang.org/gopher")
	coauthor := NewPerson("Octo Cat", "octo@github.com", "https://octodex.github.com/")
	baseURL := "https://example.com"
	feedURL := baseURL + "/feed.atom"

	entry1Date := now.Add(-72 * time.Hour)
	entry2Date := now.Add(-48 * time.Hour)
	entry3Date := now.Add(-12 * time.Hour)
	entries := []Entry{
		*NewEntry(NewEntryID(feedID, entry1Date), "Article 1", baseURL+"/blog/1", author, entry1Date, entry1Date, []string{"tech", "go"}, []byte("<em>summary</em>"), []byte("<h1>Header 1</h1>")),
		*NewEntry(NewEntryID(feedID, entry2Date), "Article 2", baseURL+"/blog/2", author, entry2Date, time.Time{}, nil, nil, []byte("<h1>Header 2</h1>")),
		*NewEntry(NewEntryID(feedID, entry3Date), "Article 3", baseURL+"/blog/3", coauthor, entry3Date, time.Time{}, []string{"dog", "cat"}, []byte("I'm a cat!"), []byte("<h1>Header 3</h1>")),
	}

	feed := NewFeed(feedID, author, title, subtitle, baseURL, feedURL, now, entries)
	if err := feed.Verify(); err != nil {
		log.Fatal(err)
	}

	out := &bytes.Buffer{}
	if err := feed.Encode(out); err != nil {
		log.Fatal(err)
	}
	got := out.String()
	want := basicBlogFeed
	if got != want {
		t.Errorf("NewFeed() returned unexpected result\n\ngot:\n%v\n\nwant:\n%v", got, want)
	}
}

const basicBlogFeed = `<?xml version="1.0" encoding="UTF-8"?>
<feed xmlns="http://www.w3.org/2005/Atom">
  <id>tag:example.com,2012-12-21:blog</id>
  <generator uri="https://github.com/denisbrodbeck/atomfeed" version="1.0">atomfeed package</generator>
  <link href="https://example.com" rel="alternate" type="text/html"></link>
  <link href="https://example.com/feed.atom" rel="self" type="application/atom+xml"></link>
  <updated>2012-12-21T08:30:15Z</updated>
  <title>example.com blog</title>
  <subtitle>Get the very latest news from the net.</subtitle>
  <author>
    <name>Go Pher</name>
    <uri>https://blog.golang.org/gopher</uri>
  </author>
  <entry>
    <id>tag:example.com,2012-12-21:blog.post-20121218083015</id>
    <title>Article 1</title>
    <link href="https://example.com/blog/1" rel="alternate" type="text/html"></link>
    <published>2012-12-18T08:30:15Z</published>
    <updated>2012-12-18T08:30:15Z</updated>
    <author>
      <name>Go Pher</name>
      <uri>https://blog.golang.org/gopher</uri>
    </author>
    <category term="tech"></category>
    <category term="go"></category>
    <summary type="html">&lt;em&gt;summary&lt;/em&gt;</summary>
    <content type="html">&lt;h1&gt;Header 1&lt;/h1&gt;</content>
  </entry>
  <entry>
    <id>tag:example.com,2012-12-21:blog.post-20121219083015</id>
    <title>Article 2</title>
    <link href="https://example.com/blog/2" rel="alternate" type="text/html"></link>
    <updated>2012-12-19T08:30:15Z</updated>
    <author>
      <name>Go Pher</name>
      <uri>https://blog.golang.org/gopher</uri>
    </author>
    <content type="html">&lt;h1&gt;Header 2&lt;/h1&gt;</content>
  </entry>
  <entry>
    <id>tag:example.com,2012-12-21:blog.post-20121220203015</id>
    <title>Article 3</title>
    <link href="https://example.com/blog/3" rel="alternate" type="text/html"></link>
    <updated>2012-12-20T20:30:15Z</updated>
    <author>
      <name>Octo Cat</name>
      <email>octo@github.com</email>
      <uri>https://octodex.github.com/</uri>
    </author>
    <category term="dog"></category>
    <category term="cat"></category>
    <summary type="html">I&#39;m a cat!</summary>
    <content type="html">&lt;h1&gt;Header 3&lt;/h1&gt;</content>
  </entry>
</feed>`

func TestNewFeedID(t *testing.T) {
	now := time.Date(2000, time.February, 26, 8, 30, 15, 0, time.UTC)

	type args struct {
		authorityName string
		creationTime  time.Time
		specific      string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "empty",
			args: args{
				authorityName: "",
				creationTime:  now,
				specific:      "",
			},
			want:    "tag:,2000-02-26:",
			wantErr: false,
		},
		{
			name: "domain",
			args: args{
				authorityName: "example.org",
				creationTime:  now,
				specific:      "blog",
			},
			want:    "tag:example.org,2000-02-26:blog",
			wantErr: false,
		},
		{
			name: "email",
			args: args{
				authorityName: "mail@example.org",
				creationTime:  now,
				specific:      "blog",
			},
			want:    "tag:mail@example.org,2000-02-26:blog",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewFeedID(tt.args.authorityName, tt.args.creationTime, tt.args.specific)
			if got.Value != tt.want {
				t.Errorf("NewFeedID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewEntryID(t *testing.T) {
	now := time.Date(2000, time.February, 26, 8, 30, 15, 0, time.UTC)

	type args struct {
		authorityName     ID
		entryCreationTime time.Time
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "empty",
			args: args{
				authorityName:     ID{},
				entryCreationTime: now,
			},
			want:    ".post-20000226083015",
			wantErr: false,
		},
		{
			name: "domain",
			args: args{
				authorityName:     NewFeedID("example.org", now, "blog"),
				entryCreationTime: now,
			},
			want:    "tag:example.org,2000-02-26:blog.post-20000226083015",
			wantErr: false,
		},
		{
			name: "email",
			args: args{
				authorityName:     NewFeedID("mail@example.org", now, "blog"),
				entryCreationTime: now,
			},
			want:    "tag:mail@example.org,2000-02-26:blog.post-20000226083015",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewEntryID(tt.args.authorityName, tt.args.entryCreationTime)
			if got.Value != tt.want {
				t.Errorf("NewEntryID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewContent(t *testing.T) {
	gif64 := `R0lGODdhAQABAIAAAP///////ywAAAAAAQABAAACAkQBADs=`
	gif, err := base64.StdEncoding.DecodeString(gif64)
	if err != nil {
		t.Error(err)
	}
	type args struct {
		contentType string
		source      string
		value       []byte
	}
	tests := []struct {
		name string
		args args
		want *Content
	}{
		{
			name: "empty",
			args: args{
				contentType: "",
				source:      "",
				value:       []byte(""),
			},
			want: nil,
		},
		{
			name: "textonly",
			args: args{
				contentType: "",
				source:      "",
				value:       []byte("some text"),
			},
			want: &Content{Value: "some text"},
		},
		{
			name: "html",
			args: args{
				contentType: "html",
				source:      "",
				value:       []byte("<h1>Header</h1>"),
			},
			want: &Content{Type: "html", Value: "<h1>Header</h1>"},
		},
		{
			name: "svg",
			args: args{
				contentType: "image/svg+xml",
				source:      "",
				value:       []byte(`<svg baseProfile="full" width="300" height="200" xmlns="http://www.w3.org/2000/svg"><rect width="100%" height="100%" fill="red"/></svg>`),
			},
			want: &Content{Type: "image/svg+xml", Value: `<svg baseProfile="full" width="300" height="200" xmlns="http://www.w3.org/2000/svg"><rect width="100%" height="100%" fill="red"/></svg>`},
		},
		{
			name: "gif",
			args: args{
				contentType: "image/gif",
				source:      "",
				value:       gif,
			},
			want: &Content{Type: "image/gif", Value: gif64, base64Encoded: true},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewContent(tt.args.contentType, tt.args.source, tt.args.value); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewContent() = %v, want %v", got, tt.want)
			}
		})
	}
}
