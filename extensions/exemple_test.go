package extensions_test

import (
	"fmt"
	"log"

	"github.com/fredbi/gflag"
	"github.com/fredbi/gflag/extensions"
	"github.com/spf13/pflag"
)

// Joint usage with pflag, for a simple bool flag
func ExampleByteSize() {
	flagVal := extensions.NewByteSizeValue(nil, 1024)
	name := "size"
	short := name[:1]
	fs := pflag.NewFlagSet("Example", pflag.ContinueOnError)
	fs.VarP(flagVal, name, short, "byte size")

	if err := fs.Parse([]string{"--size", "20GB"}); err != nil {
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

	// using GetValue()
	fmt.Println(flagVal.GetValue())

	// Output:
	// size
	// 20GB
	// size
	// 20GB
	// 20GB
	// 20000000000
}

// Example of joint usage with gflag generic type.
func Example() {
	fl := gflag.NewFlagValue(extensions.NewByteSizeValue(nil, 1024), 1024)

	name := "size"
	short := name[:1]

	fs := pflag.NewFlagSet("Example", pflag.ContinueOnError)
	pf := fs.VarPF(fl, name, short, "byte size")

	if err := fs.Parse([]string{"--size", "10GB"}); err != nil {
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
	fmt.Println(pf.Value)

	// using type's GetValue()
	fmt.Printf("%[1]T: %[1]v\n", fl.Value.GetValue()) // return uint64

	// using generic GetValue[ByteSizeValue]()
	fmt.Printf("%[1]T: %[1]v\n", uint64(fl.GetValue()))    // GetValue() returns ByteSize
	fmt.Printf("%[1]T: %[1]v\n", fl.GetValue().GetValue()) // returns uint64

	// Output:
	// size
	// 10GB
	// size
	// 10GB
	// 10GB
	// uint64: 10000000000
	// uint64: 10000000000
	// uint64: 10000000000
}
