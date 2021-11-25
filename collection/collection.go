// Package collection provides utilities for operating directly on or with collections of protobuf messages.
package collection

import "google.golang.org/protobuf/proto"

// MessageCollection encapsulates a collection of protobuf messages.
type MessageCollection []proto.Message
