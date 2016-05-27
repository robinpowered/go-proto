# Protocol Buffer Collections

[![GoDoc](https://godoc.org/github.com/mcos/protobuf-collections/pbs?status.svg)](https://godoc.org/github.com/mcos/protobuf-collections/pbs)

A `Go` library providing support for operations on collections of protocol buffer messages.

## Install
`go get -u github.com/mcos/protobuf-collections`

## Length-Prefix Frame Streaming
The initial use case provides functions for reading and writing message collections using a
length-prefixed framing technique, as discussed in the [Streaming Multiple Messages](https://developers.google.com/protocol-buffers/docs/techniques#streaming)
of the official documentation.
