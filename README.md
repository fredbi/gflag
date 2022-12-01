![Lint](https://github.com/fredbi/gflag/actions/workflows/01-golang-lint.yaml/badge.svg)
![CI](https://github.com/fredbi/gflag/actions/workflows/02-test.yaml/badge.svg)
[![Coverage Status](https://coveralls.io/repos/github/fredbi/gflag/badge.svg)](https://coveralls.io/github/fredbi/gflag)
![GitHub tag (latest by date)](https://img.shields.io/github/v/tag/fredbi/gflag)
[![Go Reference](https://pkg.go.dev/badge/github.com/fredbi/gflag.svg)](https://pkg.go.dev/github.com/fredbi/gflag)

# gflag

`pflags` with generic types.

> Yet another CLI flags library that reuses the great package `github.com/spf13/pflag`, with an interface built on go generics.
> This is not a fork, but an extension of the `pflag` functionality.

## Usage

This package is designed to be used with `pflag`. All types created by `gflag` implement the `pflag.Value` interface.

```go
	var flagVal int
	fs := pflag.NewFlagSet("", pflag.ContinueOnError)
	intFlag := gflag.NewFlagValue(&flagVal, 1)      // infer flag from underlying type int, with a default value

	fs.Var(intFlag, "integer",  "integer value")    // register the flag in pflag flagset

	_ = fs.Parse([]string{"--integer", 10})         // parse command line arguments

	fmt.Println(intFlag.GetValue())
```

You may take a look at [more examples](example_values_test.go), with slices and maps.

## What does this bring to the table?

This package provides a unified approach to strongly typed flags, using go generics.

This way, we no longer have to register CLI flags using dozen of type-specific methods.

The flag building logic is now consistent for single values, slices and maps of all types.

All primitive types (yes, complex too!) are supported.
All common types handled by `pflag` (`time.Duration`, `net.IP`, etc..) are also supported.
I have also added `time.Time`.

Slight variations in the semantics for a flag with a given underlying type may be fine-tuned with some options.

There is an `extensions` sub-package to contribute some more custom flag types. 
This one is initially populated with a `byte-size` flag. See [the example](extensions/example_test.go).

## About CLI flags

There are quite many existing CLI flag handling libraries out there.

The most popular one is (by a mile) `github.com/spf13/pflag` (originally forked from `githb.com/ogier/pflag`).

It is used by other very popular packages `spf13/viper` and `spf13/cobra`.

The second most popular (by number of imports by public repositories) is `github.com/jessevdk/go-flags`.

The approach proposed by our package is built on top `github.com/spf13/pflag` and remains interoperable with other 
great CLI-building libraries such as `viper` and `cobra`.

[Various packages and approaches to dealing with flags](docs/approaches.md)
