<h1 align="center">
  <br>
    <img src="https://github.com/denisbrodbeck/atomfeed/blob/master/logo.png?raw=true" alt="logo" width="512" height="473">
  <br>
  atomfeed – create atom syndication feeds the easy way
  <br>
  <br>
</h1>

<h4 align="center">Create Atom 1.0 compliant feeds for blogs and more.</h4>

<p align="center">
  <a href="https://godoc.org/github.com/denisbrodbeck/atomfeed"><img src="https://godoc.org/github.com/denisbrodbeck/atomfeed?status.svg" alt="GoDoc"></a>
  <a href="https://goreportcard.com/report/github.com/denisbrodbeck/atomfeed"><img src="https://goreportcard.com/badge/github.com/denisbrodbeck/atomfeed" alt="Go Report Card"></a>
</p>

## Main Features

This package:

* allows easy creation of valid Atom 1.0 feeds.
* provides convenience functions to create feeds suitable for most blogs.
* enables creation of complex atom feeds by usage of low–level structs.
* checks created feeds for most common issues (missing IDs, titles, time stamps…).
* has no external dependencies

## Installation

Import the library with

```golang
import "github.com/denisbrodbeck/atomfeed"
```

## Usage

```golang
package main

import (
	"bytes"
	"fmt"
	"log"
	"time"

	"github.com/denisbrodbeck/atomfeed"
)

func main() {
	// create constant and unique feed ID
	feedID := atomfeed.NewID("tag:example.com,2005-12-21:blog")
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
		atomfeed.NewEntryID(feedID, entry1Date), // constant and unique id
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
		atomfeed.NewEntryID(feedID, entry2Date), // constant and unique id
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
		entry1,
		entry2,
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
```

## Link

Make your atom feed discoverable by adding a `link` to your html's head:

```html
<head>
	…
	<link rel="alternate" type="application/atom+xml" title="Atom Feed" href="https://example.com/feed.atom" />
	…
</head>
```

Link to your feed from your content:

```html
<a href="https://example.com/feed.atom" type="application/atom+xml" title="Atom Feed">Subscribe with Atom</a>
```

## ID

> So what IDs should I use for my atom feed / entries?

[RFC 4287](https://tools.ietf.org/html/rfc4287#section-4.2.6) writes the following:

> When an Atom Document is relocated, migrated, syndicated, republished, exported, or imported, the content of its atom:id element MUST NOT change.

There are three requirements for an Atom ID:

1. The ID must be a valid URI (see [RFC 4287](https://tools.ietf.org/html/rfc4287#section-4.2.6)).
1. The ID must be globally unique, across all Atom feeds, everywhere, for all time.
1. The ID must never change.

> I can use permalinks, they are unique and don't change, aren't they?

Well, you *could* use permalinks as Atom IDs, but depending on how your permalinks are constructed, they *could* change. Imagine permalinks, which are automatically constructed from your `base URL` and your post's `title`. What happens, if you update the `title` of your post? The permalink to your post **changes** and thus the Atom ID of your entry changes.

> So what do I use instead? UUID? URN?

Again, you *could* use a UUID and store it along side your post, but it's not easily readable by humans. URNs require additional [registration](https://tools.ietf.org/html/rfc3406).

There is an easier way to human readable, unique and constant IDs, though: [Tag URIs](http://www.taguri.org/)

### tag URI

[Tag URIs](http://www.taguri.org/) are defined in [RFC 4151](https://tools.ietf.org/html/rfc4151).

Example-Feed: `tag:example.com,2005:blog`

Example-Post: `tag:example.com,2005:blog.post-20171224083015`

* start with `tag:`
* append an *authority name* (the domain you own or an email address): `example.com`
* append a comma `,`
* append a date, that signifies, when you had control/ownership over this *authority name*, like `2005` or `2005-02` or `2005-02-24`
* append a colon `:`
* append a specifier (anything you like): `blog`
* you've got a valid ID for an *atom:feed*: `tag:example.com,2005:blog`
  * append a dot `.`
  * append the posts creation time without special characters, turn `2017-12-24 08:30:15` into `20171224083015`
  * you've got a valid ID for an *atom:entry*: ``tag:example.com,2005:blog.post-20171224083015``

For further info check out Mark Pilgrims article on [how to make a good ID in Atom](http://web.archive.org/web/20110514113830/http://diveintomark.org/archives/2004/05/28/howto-atom-id).

## Verification

The Atom 1.0 standard defines several must–have properties of valid atom feeds
and this package allows the feed author to verify the validity of created feeds and entries.

Most common issues are (verified by this package):

* Missing or invalid `ID` on **atom:feed**
* Missing or invalid `ID` on **atom:entry**
* Missing *titles* on **atom:feed** / **atom:entry**
* Invalid time stamps (missing or not RFC3339 compliant)
* Missing *author*
* Invalid *URIs* in elements which require a valid IRI (**atom:icon**)
* Invalid *content*

```golang
package main

import (
	"log"
	"time"
	"github.com/denisbrodbeck/atomfeed"
)

func main() {
	feed := atomfeed.Feed{
		ID:      atomfeed.NewID("tag:example.com,2005-07-18:blog"),
		Title:   &atomfeed.TextConstruct{Value: "Deep Dive Into Go"},
		Updated: atomfeed.NewDate(time.Now()),
		Author:  &atomfeed.Person{Name: "Go Pher"},
	}
	// perform sanity checks on created feed
	if err := feed.Verify(); err != nil {
		log.Fatal(err)
	}
}
```

Further checks can be made with the [atom feed validator](https://validator.w3.org/feed/) from W3C. Please do run this validator, if you are constructing a complex feed.

## What is not included?

[RFC 4287](https://tools.ietf.org/html/rfc4287) defines several very advanced features, which were deliberately not implemented:

* [Securing Atom Documents](https://tools.ietf.org/html/rfc4287#section-5)
* [Extending Atom](https://tools.ietf.org/html/rfc4287#section-6)

## Credits

The Go gopher was created by [Denis Brodbeck](https://github.com/denisbrodbeck) with [gopherize.me](https://gopherize.me/), based on original artwork from [Renee French](http://reneefrench.blogspot.com/).

## License

The MIT License (MIT) — [Denis Brodbeck](https://github.com/denisbrodbeck). Please have a look at the [LICENSE](LICENSE) for more details.
