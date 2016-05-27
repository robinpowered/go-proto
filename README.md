# Protocol Buffer Collections

[![GoDoc](https://godoc.org/github.com/robinpowered/protobuf-collections?status.svg)](https://godoc.org/github.com/robinpowered/protobuf-collections)

A `Go` library providing support for operations on collections of protocol buffer messages.

## Install
`go get -u github.com/robinpowered/protobuf-collections`

## Length-Prefix Frame Streaming
The initial use case provides functions for reading and writing message collections using a
length-prefixed framing technique, as discussed in the [Streaming Multiple Messages](https://developers.google.com/protocol-buffers/docs/techniques#streaming)
of the official documentation.

## License

**protobuf-collections** is licensed under the [Apache License, Version 2.0][license-file].

--------------------------------------------------------------------------------

Copyright 2016 Robin Powered, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.




[license-file]: LICENSE
