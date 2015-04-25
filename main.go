package main

import (
	"flag"
	"fmt"
	"os"

	"git.apache.org/thrift.git/lib/go/thrift"
)

var (
	from Protocol
	to   Protocol
)

func main() {
	flag.Var(&from, "from", "protocol to convert from")
	flag.Var(&to, "to", "protocol to convert to")
	flag.Parse()

	transport := thrift.NewStreamTransport(os.Stdin, os.Stdout)

	var input, output thrift.TProtocol
	var err error

	input, err = from.Build(transport)
	if err != nil {
		panic(err)
	}
	output, err = to.Build(transport)
	if err != nil {
		panic(err)
	}

	converter := Converter{input, output}
	if err := converter.convertMessage(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
