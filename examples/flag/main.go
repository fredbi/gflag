//go:build examples

//nolint:forbidigo
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net"
	"time"

	"github.com/fredbi/gflag"
)

type cmdFlags struct {
	Verbose   bool
	Integer   int
	Counter   int
	Strings   []string
	IPs       []net.IP
	Durations map[string]time.Duration
}

var cliFlags cmdFlags

func registerFlags() {
	gfl := gflag.NewFlagValue(&cliFlags.Verbose, false)
	flag.Var(gfl, "verbose", "report verbose output")

	flag.Var(gflag.NewFlagValue(&cliFlags.Integer, 1),
		"integer", "sets a number",
	)

	// we can't do here with "flag" what we do with "pflag": a value will be required,
	// e.g. -count=1
	cfl := gflag.NewFlagValue(&cliFlags.Counter, 0, gflag.IntIsCount(true))
	flag.Var(cfl, "count", "increments a counter")

	flag.Var(gflag.NewFlagSliceValue(&cliFlags.Strings, []string{"x"}),
		"strings", "sets some strings",
	)
	flag.Var(gflag.NewFlagSliceValue(&cliFlags.IPs, []net.IP{net.ParseIP("127.0.0.1")}),
		"ips", "sets some ip addresses",
	)
	flag.Var(gflag.NewFlagMapValue(&cliFlags.Durations, map[string]time.Duration{"sec": time.Second}),
		"durations", "sets a map of durations to keys",
	)
}

/*
Usage of /tmp/go-build1248450367/b001/exe/main:

	-count
	  	increments a counter (default 0)
	-durations value
	  	sets a map of durations to keys (default [sec=1s])
	-integer
	  	sets a number (default 1)
	-ips value
	  	sets some ip addresses (default [127.0.0.1])
	-strings value
	  	sets some strings (default [x])
	-verbose
	  	report verbose output (default false)

Full example:

go run main.go -count 1 -integer 12 -count 1 -ips "8.8.8.8,2.2.2.2" -strings "a,b" -durations "hour=1h,second=1s" -durations "day=24h" -verbose
flag values parsed

	{
	 "Verbose": true,
	 "Integer": 12,
	 "Counter": 1,
	 "Strings": [
	  "a",
	  "b"
	 ],
	 "IPs": [
	  "8.8.8.8",
	  "2.2.2.2"
	 ],
	 "Durations": {
	  "day": 86400000000000,
	  "hour": 3600000000000,
	  "second": 1000000000
	 }
	}
*/
func main() {
	registerFlags()
	// pflag.CommandLine.AddGoFlagSet(flag.CommandLine)

	flag.Parse()

	fmt.Println("flag values parsed")
	bbb, _ := json.MarshalIndent(cliFlags, "", " ")
	fmt.Printf("%s\n", string(bbb))
}
