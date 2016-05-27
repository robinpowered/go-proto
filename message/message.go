// Package message provides utilities for operating directly on or with protobuf message types.
package message

import "github.com/golang/protobuf/proto"

// UnmarshalFunc is a type aliased function which is used to convert raw bytes to a protobuf message type.
// It's intended use is mainly as a callback type (when needed) when performing operations on collections of protobuf types.
// Implementing `UnmarshalFunc` is a straightforward operation.
//
//	import "io"
//	import "github.com/golang/protobuf/proto"
//	import "github.com/robinpowered/go-proto/message"
//	import "github.com/robinpowered/go-proto/stream"
//
//	// Foo is a protobuf unmarshaller
//	type Foo struct {...}
//	var f message.UnmarshalFunc = func (b []byte) (proto.Message, error){
//		var foo Foo
//		err := proto.Unmarshal(b, foo)
//
//		return foo, err
//	}
//
// The UnmarshalFunc can then be supplied to the desired function which creates an array of `Foo` types from a stream of bytes.
//
//	var r io.Reader
//	foos, err := stream.ReadLengthPrefixedCollection(r, f)
//
type UnmarshalFunc func([]byte) (proto.Message, error)
