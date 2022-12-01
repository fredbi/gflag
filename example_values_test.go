package gflag_test

import (
	"fmt"
	"log"

	"github.com/fredbi/gflag"
	"github.com/spf13/pflag"
)

// Joint usage with pflag, for a simple bool flag
func ExampleValue() {
	var flagVal bool
	name := "verbose"
	short := name[:1]
	fs := pflag.NewFlagSet("Example", pflag.ContinueOnError)

	verboseFlag := gflag.NewFlagValue(&flagVal, false)
	fl := fs.VarPF(verboseFlag, name, short, "verbose output")
	fl.NoOptDefVal = verboseFlag.GetNoOptDefVal() // allow no argument passed to the flag

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

	// using underlying value
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
