//go:build examples

//nolint:forbidigo
package main

import (
	"fmt"
	"net"
	"time"

	"github.com/fredbi/gflag"
	"github.com/spf13/pflag"
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
	pflag.Var(gfl, "verbose", "report verbose output")

	// this is still a little awkwards to set this, because pflag.Var does not support an interface to support NoOptDefVal for Value types.
	fl := pflag.Lookup("verbose")
	fl.NoOptDefVal = gfl.GetNoOptDefVal()

	pflag.Var(gflag.NewFlagValue(&cliFlags.Integer, 1),
		"integer", "sets a number",
	)

	// same clumsiness as for bool. We'd need to improve this in pflag
	cfl := gflag.NewFlagValue(&cliFlags.Counter, 0, gflag.IntIsCount(true))
	pflag.Var(cfl, "count", "increments a counter")
	fl = pflag.Lookup("count")
	fl.NoOptDefVal = cfl.GetNoOptDefVal()

	pflag.Var(gflag.NewFlagSliceValue(&cliFlags.Strings, []string{"x"}),
		"strings", "sets some strings",
	)
	pflag.Var(gflag.NewFlagSliceValue(&cliFlags.IPs, []net.IP{net.ParseIP("127.0.0.1")}),
		"ips", "sets some ip addresses",
	)
	pflag.Var(gflag.NewFlagMapValue(&cliFlags.Durations, map[string]time.Duration{"sec": time.Second}),
		"durations", "sets a map of durations to keys",
	)
}

/*
go run main.go --help
Usage of /tmp/go-build3134170569/b001/exe/main:

	--count count                  increments a counter (default 0)
	--durations stringToDuration   sets a map of durations to keys (default [sec=1s])
	--integer int                  sets a number (default 1)
	--ips ipSlice                  sets some ip addresses (default [127.0.0.1])
	--strings strings              sets some strings (default [x])
	--verbose                      report verbose output

pflag: help requested
exit status 2

Full example:

	go run main.go --count --integer 12 --count \
	  --ips 8.8.8.8,2.2.2.2 --strings a,b \
	  --durations hour=1h,second=1s --durations day=24h
*/
func main() {
	registerFlags()

	pflag.Parse()

	fmt.Println("flag values parsed")
	fmt.Printf("%#v\n", cliFlags)
}
