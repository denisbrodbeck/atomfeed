package atomfeed_test

import (
	"bytes"
	"fmt"
	"log"
	"time"

	"github.com/denisbrodbeck/atomfeed"
)

// Create a basic atom syndication feed suitable for blogs:
func Example() {
	// create constant and unique feed ID
	authority := "example.com"
	owned := time.Date(2005, time.December, 21, 8, 30, 15, 0, time.UTC)
	feedID := atomfeed.NewFeedID(authority, owned, "blog")
	// set basic feed properties
	title := "example.com blog"
	subtitle := "Get the very latest news from the net."
	author := atomfeed.NewPerson("Go Pher", "", "https://blog.golang.org/gopher")
	coauthor := atomfeed.NewPerson("Octo Cat", "octo@github.com", "https://octodex.github.com/")
	updated := time.Date(2015, time.March, 21, 8, 30, 15, 0, time.UTC)
	baseURL := "https://example.com"
	feedURL := baseURL + "/feed.atom"

	entry1Date := time.Date(2012, time.October, 21, 8, 30, 15, 0, time.UTC)
	entry1 := atomfeed.NewEntry(
		atomfeed.NewEntryID(*feedID, entry1Date), // constant and unique id
		"Article 1",                              // title
		baseURL+"/post/1",                        // permalink
		author,                                   // author of the entry/post
		entry1Date,                               // updated date – mandatory
		entry1Date,                               // published date – optional
		[]string{"tech", "go"},                   // categories
		[]byte("<em>go go go</em>"),              // summary – optional
		[]byte("<h1>Header 1</h1>"),              // content
	)
	entry2Date := time.Date(2012, time.December, 21, 8, 30, 15, 0, time.UTC)
	entry2 := atomfeed.NewEntry(
		atomfeed.NewEntryID(*feedID, entry2Date), // constant and unique id
		"Article 2",                              // title
		baseURL+"/post/2",                        // permalink
		coauthor,                                 // author of the entry/post
		entry2Date,                               // updated date – mandatory
		entry2Date,                               // published date – optional
		[]string{"cat", "dog"},                   // categories – optional
		[]byte("I'm a cat!"),                     // summary – optional
		[]byte("<h1>Header 2</h1>"),              // content
	)
	entries := []atomfeed.Entry{
		*entry1,
		*entry2,
	}

	feed := atomfeed.NewFeed(feedID, author, title, subtitle, baseURL, feedURL, updated, entries)
	// most atom elements support language attributes (optional)
	feed.CommonAttributes = &atomfeed.CommonAttributes{Lang: "en"}
	// perform sanity checks on created feed
	if err := feed.Verify(); err != nil {
		log.Fatal(err)
	}

	out := &bytes.Buffer{}
	// serialize XML data into stream
	if err := feed.Encode(out); err != nil {
		log.Fatal(err)
	}
	fmt.Print(out.String())
}

// Create new constant and unique feed id:
func ExampleNewFeedID() {
	authorityName := "example.com"
	ownedAt := time.Date(2005, time.July, 18, 0, 0, 0, 0, time.UTC)
	specifier := "blog"
	id := atomfeed.NewFeedID(authorityName, ownedAt, specifier)
	fmt.Println(id.Value)
	// Output: tag:example.com,2005-07-18:blog
}

// Create new constant and unique entry id:
func ExampleNewEntryID() {
	feedID := atomfeed.ID{Value: "tag:example.com,2005-07-18:blog"}
	entryCreationTime := time.Date(2017, time.December, 21, 8, 30, 15, 0, time.UTC)
	id := atomfeed.NewEntryID(feedID, entryCreationTime)
	fmt.Println(id.Value)
	// Output: tag:example.com,2005-07-18:blog.post-20171221083015
}

// Add attributes like "lang" to feed or entry elements:
func ExampleCommonAttributes() {
	feed := atomfeed.Feed{
		ID:      &atomfeed.ID{Value: "tag:example.com,2005-07-18:blog"},
		Title:   &atomfeed.TextConstruct{Value: "Deep Dive Into Go"},
		Updated: atomfeed.NewDate(time.Now()),
		Author:  &atomfeed.Person{Name: "Go Pher"},
	}
	// add language attribute
	feed.CommonAttributes = &atomfeed.CommonAttributes{Lang: "en"}
}

// Create new category / tag:
func ExampleNewCategory() {
	category := atomfeed.NewCategory("golang")
	fmt.Println(category.Term)
	// Output: golang
}

// Create RFC 3339 compliant date:
func ExampleNewDate() {
	now := time.Date(2017, time.December, 22, 8, 30, 15, 0, time.UTC)
	date := atomfeed.NewDate(now)
	fmt.Println(date.Value)
	// Output: 2017-12-22T08:30:15Z
}
