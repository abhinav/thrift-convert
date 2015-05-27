package main

import (
	"flag"
	"os"

	"github.com/apache/thrift/lib/go/thrift"
)

var (
	from     Protocol
	to       Protocol
	isStruct = flag.Bool("struct", false, "convert a struct only")
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
	if *isStruct {
		if err := converter.convertType(thrift.STRUCT); err != nil {
			panic(err)
		}
	} else {
		if err := converter.convertMessage(); err != nil {
			panic(err)
		}
	}

	if err := output.Flush(); err != nil {
		panic(err)
	}
}
