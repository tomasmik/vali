// Package vali is a simple and easily extendable struct validation tool.
// It only has the basic validation predefined and the package user
// is expected to define other validation funcs himself.
//
// It's meant to serve as a sort of scripting language in your
// struct field tags, allowing you to define a path on how to
// validate a field.
//
// The user is expected to define the new validator instance
// and register all of his validation funcs before actually
// using it, as the maps in which they're held are not thread safe.
// That can and should be done in the packages `init()` func.
//
// Most of the documentation can be found by reading the source code,
// for example the documentation on what tags there are and how to use them
// can be found in the `tagfn.go` file, while the information on what
// the tags are made of can be found in `vali.go`.
//
// Key points to help you use the validator:
// Tags are validated in order.
// Tags are seperated by `|`
// Tag values by `,`
// There are special tags: `*` to point to another struct field and `>` to validate slice elements
//
// Example tag: `vali:"min=2|>|one_of=a,b"`. It will validate that a given slice is
// longer than 2 elements and that it's values are either `a` or `b` strings.
package vali
