// Package atomfeed creates atom syndication feeds (Atom 1.0 RFC4287).
//
// https://github.com/denisbrodbeck/atomfeed
//
// https://tools.ietf.org/html/rfc4287
//
// This package allows easy creation of valid Atom 1.0 feeds.
// It provides functions to create feeds suitable for most blogs.
// Direct usage of low–level structs allows the creation of more complex atom feeds.
//
// The Atom 1.0 standard defines several must–have properties of valid atom feeds
// and this package allows the feed author to verify the validity of created feeds and entries
// and to check for most common issues (missing IDs, titles, timestamps…).
// Further validation with the final output of this package
// can be done with the online validator https://validator.w3.org/feed/ from W3C.
package atomfeed // import "github.com/denisbrodbeck/atomfeed"
