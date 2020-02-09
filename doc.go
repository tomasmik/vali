// Package vali is a simple and easily extendable struct validation package.
// It only has the basic validation predefined and the package user
// is expected to define other validation funcs himself.
//
// The user is expected to define the new validator instance
// and register all of his validation funcs before actually
// using it as the maps in which they're held are not thread safe.
// That can be and should be done in the packages `init()` func.
//
// Documentation on premade tags can be found in the `tagfn.go` file.
// Documentation on adding validators and tag can be found in `vali.go`.
package vali
