//go:build examples

//nolint:forbidigo
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"github.com/fredbi/gflag"
	"github.com/fredbi/gflag/extensions"
	"github.com/jessevdk/go-flags"
)

type cmdFlags struct {
	Verbose     *gflag.Value[bool]             `long:"verbose" description:"verbose output"`
	Integer     *gflag.Value[int]              `long:"integer" description:"sets a number"`
	Counter     *gflag.Value[int]              `long:"count" description:"increments a counter"`
	Strings     *gflag.SliceValue[string]      `long:"strings" description:"sets some strings"`
	IPs         *gflag.SliceValue[net.IP]      `long:"ips" description:"sets some ip addresses"`
	Durations   *gflag.MapValue[time.Duration] `long:"durations" description:"sets a map of durations to keys"`
	WithDefault *extensions.ByteSizeValue      `long:"with-default" description:"default in struct tag" default:"1GB"`
	WithChoice  *gflag.Value[string]           `long:"with-choice" description:"default in struct tag" choice:"json" choice:"yaml"`
}

func registerFlags() *cmdFlags {
	return &cmdFlags{
		// default values handled programmatically
		Verbose:   gflag.NewFlagValue(nil, false),
		Integer:   gflag.NewFlagValue(nil, 1),
		Counter:   gflag.NewFlagValue(nil, 0, gflag.IntIsCount(true)), // TODO: does not work
		Strings:   gflag.NewFlagSliceValue(nil, []string{"x"}),
		IPs:       gflag.NewFlagSliceValue(nil, []net.IP{net.ParseIP("127.0.0.1")}),
		Durations: gflag.NewFlagMapValue(nil, map[string]time.Duration{"sec": time.Second}),

		// Default values handled from struct tag do not work with the generic type extension: go-flags rely on reflection and the check doesn't include structs.
		// However, it is okay to use extensions that are backed by some primitive type, slice or map.
		WithDefault: extensions.NewByteSizeValue(nil, 0),

		// List of choices work fine
		WithChoice:  gflag.NewFlagValue(nil, "json"),
	}
}

/*
go run main.go --help
Usage:
  main [OPTIONS]

Application Options:
      --verbose=                verbose output (default: false)
      --integer=                sets a number (default: 1)
      --count=                  increments a counter (default: 0)
      --strings=                sets some strings (default: [x])
      --ips=                    sets some ip addresses (default: [127.0.0.1])
      --durations=              sets a map of durations to keys (default: [sec=1s])
      --with-default=           default in struct tag (default: 1GB)
      --with-choice=[json|yaml] default in struct tag (default: json)

Help Options:
  -h, --help                    Show this help message


 go run main.go --count=1 --integer 12 --count=1   --ips 8.8.8.8,2.2.2.2 --strings a,b   --durations hour=1h,second=1s --durations day=24h --with-default=5MB --with-choice=yaml
flag values parsed
{
 "Verbose": "false",
 "Integer": "12",
 "Counter": "1",
 "Strings": "[a,b]",
 "IPs": "[8.8.8.8,2.2.2.2]",
 "Durations": "[day=24h0m0s,hour=1h0m0s,second=1s]",
 "WithDefault": "5MB",
 "WithChoice": "yaml"
}
*/
func main() {
	cliFlags := registerFlags()

	parser := flags.NewParser(cliFlags, flags.Default)
	if _, err := parser.Parse(); err != nil {
		os.Exit(1)
	}

	fmt.Println("flag values parsed")
	bbb,err := json.MarshalIndent(cliFlags, "", " ")
	if err != nil {
		log.Fatalf("%v", err)
	}

	fmt.Printf("%s\n", string(bbb))
}
