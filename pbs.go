// Package pbs provides utility methods and operations with may be performed on collections, or data representing collections (such as a binary stream) of protocol buffers.
package pbs

import (
	"encoding/binary"
	"io"

	"github.com/golang/protobuf/proto"
)

const (
	// uint32ByteSize is the size, in bytes of an unsigned 32-bit integer.
	uint32ByteSize = 4
)

// writeLengthPrefixedMessage writes a message to the supplied writer. It prefixes the message with a four byte little-endian representation of the size of the message. The total number of bytes written is returned, along with any error arising.
func writeLengthPrefixedMessage(w io.Writer, pb proto.Message) (int, error) {
	s := proto.Size(pb)
	n := uint32(s)

	err := binary.Write(w, binary.LittleEndian, n)
	if nil != err {
		return 0, err
	}

	b, err := proto.Marshal(pb)
	if nil != err {
		return 0, err
	}

	_, err = w.Write(b)
	if nil != err {
		return 0, err
	}

	return s + uint32ByteSize, nil
}

// WriteLengthPrefixedCollection writes the collection of protobuf messages to the supplied writer.
// The write operation uses length-prefixed framing. This means that each protocol buffer message is prefixed
// by its length. This implementation encodes the length as a four-byte little-endian field, representing an unsigned 32 bit integer.
// The total number of bytes (including prefixes) written to the buffer is returned, along with any error arising.
func WriteLengthPrefixedCollection(w io.Writer, pbs []proto.Message) (n int, err error) {
	for _, pb := range pbs {
		i, err := writeLengthPrefixedMessage(w, pb)
		if nil != err {
			return n, err
		}

		n += i
	}

	return n, nil
}

// UnmarshalFunc is a type aliased function which is used to convert raw bytes to a protobuf message type.
// It's intended use is mainly as a callback type (when needed) when performing operations on collections of protobuf types.
type UnmarshalFunc func([]byte) (proto.Message, error)

// ReadLengthPrefixedCollection reads a collection of protocol buffer messages from the supplied reader.
// Each message is presumed prefixed by a 4 byte little-endian field (an unsigned 32 bit integer) which represents the size of the ensuing message.
// The UnmarshalFunc argument is a supplied callback used to convert the raw bytes read as a message to the desired message type.
// The protocol buffer message collection is returned, along with any error arising.
func ReadLengthPrefixedCollection(r io.Reader, f UnmarshalFunc) (pbs []proto.Message, err error) {
	for {
		var s uint32
		err := binary.Read(r, binary.LittleEndian, &s)
		if io.EOF == err {
			return pbs, nil
		} else if nil != err {
			return nil, err
		}

		b := make([]byte, s)

		_, err = io.ReadFull(r, b)
		if io.EOF == err {
			return pbs, nil
		} else if nil != err {
			return nil, err
		}

		pb, err := f(b)
		if nil != err {
			return nil, err
		}

		pbs = append(pbs, pb)
	}
}
