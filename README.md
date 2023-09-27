![Lint](https://github.com/fredbi/gflag/actions/workflows/01-golang-lint.yaml/badge.svg)
![CI](https://github.com/fredbi/gflag/actions/workflows/02-test.yaml/badge.svg)
[![Coverage Status](https://coveralls.io/repos/github/fredbi/gflag/badge.svg)](https://coveralls.io/github/fredbi/gflag)
![Vulnerability Check](https://github.com/fredbi/gflag/actions/workflows/03-govulncheck.yaml/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/fredbi/gflag)](https://goreportcard.com/report/github.com/fredbi/gflag)

![GitHub tag (latest by date)](https://img.shields.io/github/v/tag/fredbi/gflag)
[![Go Reference](https://pkg.go.dev/badge/github.com/fredbi/gflag.svg)](https://pkg.go.dev/github.com/fredbi/gflag)
[![license](http://img.shields.io/badge/license/License-Apache-yellow.svg)](https://raw.githubusercontent.com/fredbi/go-cli/master/LICENSE.md)

# gflag

`pflags` with generic types.

> OK so this is yet another CLI flags library...
>
> We do not reinvent the wheel: this module reuses the great package `github.com/spf13/pflag`,
> with an interface built on go generics.
>
> The main idea is to simplify the `pflag` interface, with less things to remember about flag types.
>
> So this is not a fork or drop-in replacement, but rather an extension of the `pflag` functionality.

## Usage

This package is designed to be used together with `github.com/spf13/pflag`.
All types created by `gflag` implement the `pflag.Value` interface.

Use-case: build idiomatic `cobra` CLIs with a simpler declaration of flags.
See how [go-cli](https://github.com/fredbi/go-cli) achieves this.

```go
import (
	"fmt"

	"github.com/fredbi/gflag"
	"github.com/spf13/pflag"
)

var flagVal int

fs := pflag.NewFlagSet("", pflag.ContinueOnError)
intFlag := gflag.NewFlagValue(&flagVal, 1)      // infer flag from underlying type int, with a default value

fs.Var(intFlag, "integer",  "integer value")    // register the flag in pflag flagset

_ = fs.Parse([]string{"--integer", 10})         // parse command line arguments

fmt.Println(intFlag.GetValue())                 // the flag knows the type of the value
```

With `pflag` this piece of code would look very much similar.

```go
import (
	"fmt"

	"github.com/spf13/pflag"
)

var flagVal int

fs := pflag.NewFlagSet("", pflag.ContinueOnError)

fs.IntVar(&flagVal, "integer",  "integer value") // register the flag in pflag flagset

_ = fs.Parse([]string{"--integer", 10})         // parse command line arguments

fmt.Println(fs.GetInt("integer"))               // has to know an integer type is expected
```

You may take a look at [more examples](example_values_test.go), with slices and maps.

## What does this bring to the table?

This package provides a unified approach to strongly typed flags, using go generics.

This way, we no longer have to register CLI flags using dozen of type-specific methods, just three:
* `NewFlagValue()`: for scalar values
* `NewFlagSliceValue()`: for flags as lists
* `NewFlagMapValue()`: for flags as map (e.g. key=value pairs)

The flag building logic is consistent for single values, slices and maps of all types.

* **All primitive types** (yes, complex too!) are supported.
* All common types handled by `pflag` (`time.Duration`, `net.IP`, etc..) are also supported.
* I have also added `time.Time`
* Support any extension built for `pflag` (arbitrary types with the `pflag.Value` interface)


Variations in the semantics for a flag with a given underlying type may be fine-tuned with options.


There is an `extensions` sub-package to contribute some custom extra flag types. 
It is now populated with a `byte-size` flag. See [the example](extensions/example_test.go).

## [About CLI flags](./docs/about.md)
