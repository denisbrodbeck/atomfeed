package atomfeed_test

import (
	"bytes"
	"fmt"
	"log"
	"time"

	"github.com/denisbrodbeck/atomfeed"
)

func Example() {
	authority := "example.com"
	now := time.Date(2012, time.December, 21, 8, 30, 15, 0, time.UTC)
	feedID := atomfeed.NewFeedID(authority, now, "blog")
	title := "example.com blog"
	subtitle := "Get the very latest news from the net."
	author := atomfeed.NewPerson("Go Pher", "", "https://blog.golang.org/gopher")
	baseURL := "https://example.com"
	feedURL := baseURL + "/feed.atom"

	feed := atomfeed.NewFeed(feedID, author, title, subtitle, baseURL, feedURL, now, nil)
	// perform some sanity checks
	if err := feed.Verify(); err != nil {
		log.Fatal(err)
	}

	out := &bytes.Buffer{}
	if err := feed.Encode(out); err != nil {
		log.Fatal(err)
	}
	fmt.Print(out.String())
}
