#!/bin/sh

for os in linux darwin; do
	echo "Building for $os"
	GOOS="$os" GOARCH=amd64 go build -o build/"$os"/thrift-convert

	tar -cf build/thrift-convert-"$os".tar -C build/"$os"/ thrift-convert
	bzip2 build/thrift-convert-"$os".tar 
done
