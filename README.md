# vali

Vali is yet another validation package which strives to be different from others.
Main idea of this validator is that it's small, easily extendable and doesn't bring
an excessive amount of validators that most people do not need. Instead it
encourages the user to add validator that he himself needs if the default ones are not enough.

## Install

With a [correctly configured](https://golang.org/doc/install#testing) Go toolchain:

```sh
go get github.com/tomasmik/vali
```

## Features

Currently Vali comes with premade validation tags:
* required
* required_without
* optional
* min
* max
* one_of
* eq
* neq

Multiple validation tags can be added for a single struct field
and you're allow to validate against given values or other struct fields.

## Basic usage

As mentioned the package allows the user to define tags which
behave as defined by the package itself or the user.
The tags to do have a syntax which is essential when using them.

* Fields with no `vali` tag or with `vali:"-"` will be ignored
* Seperating validators can be done with the `|` symbol - `vali:"required|min=1|max=5"`
* Pointing to other struct fields can be do by using the `*` symbol - `vali:"required_without=*Foo"`
* Seperating validator values can be done by using the `,` symbol - `vali:"required|one_of=1,2,3"`

**Fields must be exported (or else they're ignored) and validate method only accepts pointers to structs**

#### Tag Example 1

Validate that `First` is more than 2, but ignore it if it's nil:

```go
	type foo struct {
		First *int `vali:"optional|max=2"`
	}
```

#### Tag Example 2

Validate that `Second` is less than First (note that it returns `false` if First is nil):

```go
	type foo struct {
		First  *int
		Second int `vali:"max=*First"`
	}
```

## Current state

The package is in it's infancy while it should work
there are some missing parts that will be added:

* General clean up of `types.go` and `tagfn.go`
* More tests
* Better documenation
* Custom errors
* `Dive` function to validate the inside of maps/structs using tags

## Contributing
If you want to improve any part of this validator you're free to create a Pull Request or an Issue.
