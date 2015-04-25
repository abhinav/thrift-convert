A very simple tool to convert payloads between different Thrift protocols.

Usage:

    thrift-convert -from PROTOCOL -to PROTOCOL

Where `PROTOCOL` is one of `binary`, `compact`, and `json`. This will read from
stdin and write to stdout.

Example:

    cat binary_payload.txt | thrift-convert -from binary -to json
