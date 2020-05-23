# vali
[![Built with Spacemacs](https://cdn.rawgit.com/syl20bnr/spacemacs/442d025779da2f62fc86c2082703697714db6514/assets/spacemacs-badge.svg)](http://spacemacs.org)

Vali is yet another validation package which strives to be different from others.
Main idea of this validator is that it's small, easily extendable and doesn't bring
an excessive amount of validators that most people do not need. Instead it
encourages the user to add validator that he himself needs if the default ones are not enough.

For documentation explore the tests, comments and [godoc](https://pkg.go.dev/github.com/tomasmik/vali?tab=doc)

## Install

With a [correctly configured](https://golang.org/doc/install#testing) Go toolchain:

```sh
go get github.com/tomasmik/vali
```

## Features

Currently Vali comes with premade validation tags:
* required (validate that value if not nil/default)
* required_without (validate that value not nil/default if *Val is)
* optional (skip if nil/default)
* min (validate that value is at least n)
* max (validate that value is below n)
* one_of (validate that value is one of given)
* none_of (validate that value is none of given)
* eq (validate that value is equal)
* neq (validate that value is not equal)
* dups (validate for duplicates in a slice)

Tags behavior:
* Multiple validation tags can be added for a single struct field.
* Validation tags are applied in order so you can chain them however you like.

Errors:
* To return a custom error, you can use the defined `BubbleErr` function.
* To force skip validation for a certain field, you can now return `ErrSkipFurther`.

Special tags:
* `>` - Allows you to validate the contents of a slice/array.
* `*` - Allows you to point to another struct field to validate against or with it.

## Basic usage

You can create a new validator with the `func New()` which will also pre-seed it with default
validation funcs or with `func NewEmpty()` which will return an empty validator instance
allowing you to define the validation workflow yourself.

As mentioned the package allows the user to define tags which
behave as defined by the package itself or the user.
The tags do have a syntax which is essential when using them.

* Fields with no `vali` tag or with `vali:"-"` will be ignored (user can change `vali` to any tag he wants)
* Seperating validators can be done with the `|` symbol - `vali:"required|min=1|max=5"`
* Pointing to other struct fields can be do by using the `*` symbol - `vali:"required_without=*Foo"`
* Seperating validator values can be done by using the `,` symbol - `vali:"required|one_of=1,2,3"`
* Validating slice and array elements is also possible (though a little experimental) by adding `>` to the validation tag - `vali:">|one_of=1,2"`

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

#### Tag Example 3

Validate that a slice length is more than the `min` amount of elements and all of them fit the `one_of` tag

```go
	type foo struct {
		First []string `vali:"min=2|>|one_of=a,b"`
	}
```

#### Utils 

`utils.go` file exposes certain util functions that are used in the package itself
and can be used by the package user if any shortcuts are needed for type conversion.

List of util functions which are exposed:

```go
func GetInt(s interface{}) (int64, bool)
func GetUInt(s interface{}) (uint64, bool) 
func GetUIntFallback(s interface{}) (uint64, bool)
func GetFloat(s interface{}) (float64, bool)
func GetString(s interface{}) string 
func DerefInterface(s interface{}) interface{} 
```

Documentation for these functions can be found in `utils.go` file.

## Current state

Other than general clean up and making the whole
validation process smoother with less allocations and conversions
time should be spent on:

* Adding more/better tests.
* Better documenation.
* Improving the error aggregation if possible.

## Contributing
If you want to improve any part of this validator you're free to create a Pull Request or an Issue.
