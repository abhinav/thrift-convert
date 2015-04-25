package main

import (
	"fmt"
	"strings"

	"git.apache.org/thrift.git/lib/go/thrift"
)

// Protocol represents a Thrift protocol type
type Protocol int

const (
	// Binary represents TBinaryProtocol
	Binary Protocol = iota
	// Compact represents TCompactProtocol
	Compact
	// JSON represents TJSONProtocol
	JSON
)

//go:generate stringer -type Protocol

// Set the value of the given Protocol based on the given string name.
func (p *Protocol) Set(s string) error {
	switch strings.ToLower(s) {
	case "binary":
		*p = Binary
	case "compact":
		*p = Compact
	case "json":
		*p = JSON
	default:
		return fmt.Errorf("unknown protocol '%s'", s)
	}
	return nil
}

// Build the TProtocol represented by this Protocol using the given Thrift
// transport.
func (p Protocol) Build(t thrift.TTransport) (thrift.TProtocol, error) {
	switch p {
	case Binary:
		return thrift.NewTBinaryProtocolTransport(t), nil
	case Compact:
		return thrift.NewTCompactProtocol(t), nil
	case JSON:
		return thrift.NewTJSONProtocol(t), nil
	default:
		return nil, fmt.Errorf("unsupported protocol '%q'", p)
	}
}
