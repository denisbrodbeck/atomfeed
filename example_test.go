package atomfeed_test

import (
	"bytes"
	"fmt"
	"log"
	"time"

	"github.com/denisbrodbeck/atomfeed"
)

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
		&entry1Date,                              // updated date – mandatory
		&entry1Date,                              // published date – optional
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
		&entry2Date,                              // updated date – mandatory
		&entry2Date,                              // published date – optional
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
