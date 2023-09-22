package gflag_test

import (
	"fmt"
	"log"

	"github.com/fredbi/gflag"
	"github.com/spf13/pflag"
)

// Joint usage with pflag, for a simple bool flag
func ExampleValue() {
	// variable to store the value of the flag
	var flagVal bool

	// Custom pflag.FlagSet.
	// Simple CLIs may just use the default pre-baked package-level FlagSet for the command line.
	name := "verbose"
	short := name[:1]
	fs := pflag.NewFlagSet("Example", pflag.ContinueOnError)

	// declare a new generic flag: type is inferred from the provided default value
	verboseFlag := gflag.NewFlagValue(&flagVal, false)

	// register this flag into the FlagSet
	fl := fs.VarPF(verboseFlag, name, short, "verbose output")
	fl.NoOptDefVal = verboseFlag.GetNoOptDefVal() // allow no argument passed to the flag

	// parse args from the command line.
	// Simple CLIs may just use the default, with pflag.Parse() from the command line arguments.
	if err := fs.Parse([]string{"--verbose"}); err != nil {
		log.Fatalln("parsing error:", err)
	}

	// retrieve parsed values
	flag := fs.Lookup(name)
	fmt.Println(flag.Name)
	fmt.Println(flag.Value)

	flag = fs.ShorthandLookup(short)
	fmt.Println(flag.Name)
	fmt.Println(flag.Value)

	// various ways to retrieve the parsed value

	// using the variable used for storing the value
	fmt.Println(flagVal)

	// using GetValue[bool]()
	fmt.Println(verboseFlag.GetValue())

	// Output:
	// verbose
	// true
	// verbose
	// true
	// true
	// true
}

// Simple int flag
func ExampleAddValueFlag() {
	const name = "integer"
	short := name[:1]
	const usage = "an integer flag"

	fs := pflag.NewFlagSet("Example", pflag.ContinueOnError)

	// add flag without a preallocated variable: all interaction is performed via the flag
	gfl, fl := gflag.AddValueFlag(
		fs,
		gflag.NewFlagValue(nil, 5), // create a generic flag of type integer
		name, short, usage,         // the usual specification for the flag
	)

	if err := fs.Parse([]string{"--integer", "12"}); err != nil {
		log.Fatalln("parsing error:", err)
	}

	// retrieve parsed values from name
	flag := fs.Lookup(name)
	fmt.Println(flag.Name)
	fmt.Println(flag.Value)

	// retrieve parsed value from short name
	flag = fs.ShorthandLookup(short)
	fmt.Println(flag.Name)
	fmt.Println(flag.Value)

	// various ways to retrieve the parsed value

	// using underlying value
	fmt.Println(fl.Value)

	// using FlagSet.GetInt() (old way)
	val, err := fs.GetInt(name)
	if err != nil {
		log.Fatalln("flag type error:", err)
	}
	fmt.Printf("%v (%T)\n", val, val)

	// using GetValue[int]() (new way)
	val2 := gfl.GetValue()
	fmt.Printf("%v (%T)\n", val2, val2)

	// Output:
	// integer
	// 12
	// integer
	// 12
	// 12
	// 12 (int)
	// 12 (int)
}

// Joint usage with pflag, for a string array flag
func ExampleSliceValue() {
	var flagVal []string

	name := "strings"
	short := name[:1]
	fs := pflag.NewFlagSet("Example", pflag.ContinueOnError)

	stringsFlag := gflag.NewFlagSliceValue(&flagVal, []string{"a", "b"})
	fs.VarP(stringsFlag, name, short, "csv input strings")

	if err := fs.Parse([]string{"--strings", "d,e,f", "--strings", "g,h"}); err != nil {
		log.Fatalln("parsing error:", err)
	}

	// retrieve parsed values
	flag := fs.Lookup(name)
	fmt.Println(flag.Name)
	fmt.Println(flag.Value)

	flag = fs.ShorthandLookup(short)
	fmt.Println(flag.Name)
	fmt.Println(flag.Value)

	// using underlying value
	fmt.Println(flagVal)

	// using GetValue[bool]()
	fmt.Println(stringsFlag.GetValue())

	// Output:
	// strings
	// [d,e,f,g,h]
	// strings
	// [d,e,f,g,h]
	// [d e f g h]
	// [d e f g h]
}

// Joint usage with pflag, for a key=value map of integers
func ExampleMapValue() {
	var flagVal map[string]int

	name := "map"
	short := name[:1]
	fs := pflag.NewFlagSet("Example", pflag.ContinueOnError)

	mapFlag := gflag.NewFlagMapValue(&flagVal, map[string]int{"a": 5, "b": 2})
	fs.VarP(mapFlag, name, short, "map of integers")

	if err := fs.Parse([]string{"--map", "f=1,g=4", "--map", "e=3"}); err != nil {
		log.Fatalln("parsing error:", err)
	}

	// retrieve parsed values
	flag := fs.Lookup(name)
	fmt.Println(flag.Name)
	fmt.Println(flag.Value)

	flag = fs.ShorthandLookup(short)
	fmt.Println(flag.Name)
	fmt.Println(flag.Value)

	// using underlying value
	fmt.Println(flagVal)

	// using GetValue[bool]()
	fmt.Println(mapFlag.GetValue())

	// Output:
	// map
	// [e=3,f=1,g=4]
	// map
	// [e=3,f=1,g=4]
	// map[e:3 f:1 g:4]
	// map[e:3 f:1 g:4]

}

// Example of using gflag options.
func ExampleOption() {
	// example with int used with count semantics
	var flagVal int

	name := "count"
	short := name[:1]
	fs := pflag.NewFlagSet("Example", pflag.ContinueOnError)

	countFlag := gflag.NewFlagValue(&flagVal, 0, gflag.IntIsCount(true))
	fl := fs.VarPF(countFlag, name, short, "increment counter")
	fl.NoOptDefVal = countFlag.GetNoOptDefVal() // allow no argument passed to the flag

	if err := fs.Parse([]string{"--count", "--count", "--count"}); err != nil {
		log.Fatalln("parsing error:", err)
	}

	// retrieve parsed values
	flag := fs.Lookup(name)
	fmt.Println(flag.Name)
	fmt.Println(flag.Value)

	flag = fs.ShorthandLookup(short)
	fmt.Println(flag.Name)
	fmt.Println(flag.Value)

	// using underlying value
	fmt.Println(flagVal)

	// using GetValue[bool]()
	fmt.Println(countFlag.GetValue())

	// Output:
	// count
	// 3
	// count
	// 3
	// 3
	// 3
}
