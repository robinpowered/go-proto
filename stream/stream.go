package stream

import (
	"encoding/binary"
	"io"

	"github.com/golang/protobuf/proto"
	"github.com/robinpowered/go-proto/collection"
	"github.com/robinpowered/go-proto/message"
)

// messageLengthPrefixByteSize is the size, in bytes, of the length prefix to use when employing length-prefix message framing. This prefix is an
// unsigned 32 bit integer, so it's always 4.
const messageLengthPrefixByteSize = 4

// LengthPrefixFramedCollectionSize calculates the size of the expected stream produced when encoding this collection using length prefixed framing.
// The implementation assumes a four byte length prefix field, representing an unsigned 32 bit integer.
func LengthPrefixFramedCollectionSize(mc collection.MessageCollection) (n int) {
	for _, m := range mc {
		n += LengthPrefixFramedMessageSize(m)
	}

	return n
}

// LengthPrefixFramedMessageSize calculates the expected total size of a length-prefixed framed protobuf message within a stream.
// It performs a simple addition of the length of the prefix field and the protobuf message size.
// The implementation assumes a four byte length prefix field, representing an unsigned 32 bit integer.
func LengthPrefixFramedMessageSize(m proto.Message) int {
	return messageLengthPrefixByteSize + proto.Size(m)
}

// ReadLengthPrefixedCollection reads a collection of protocol buffer messages from the supplied reader.
// Each message is presumed prefixed by a 4 byte little-endian field (an unsigned 32 bit integer) which represents the size of the ensuing message.
// The UnmarshalFunc argument is a supplied callback used to convert the raw bytes read as a message to the desired message type.
// The protocol buffer message collection is returned, along with any error arising.
// For more detailed information on this approach, see the official protocol buffer documentation https://developers.google.com/protocol-buffers/docs/techniques#streaming.
func ReadLengthPrefixedCollection(r io.Reader, f message.UnmarshalFunc) (pbs collection.MessageCollection, err error) {
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

// WriteLengthPrefixedCollection writes the collection of protobuf messages to the supplied writer.
// The write operation uses length-prefixed framing. This means that each protocol buffer message is prefixed
// by its length. This implementation encodes the length as a four byte little-endian field, representing an unsigned 32 bit integer.
// The total number of bytes (including prefixes) written to the buffer is returned, along with any error arising.
// For more detailed information on this approach, see the official protocol buffer documentation https://developers.google.com/protocol-buffers/docs/techniques#streaming.
func WriteLengthPrefixedCollection(w io.Writer, pbs collection.MessageCollection) (n int, err error) {
	for _, pb := range pbs {
		i, err := WriteLengthPrefixedMessage(w, pb)
		if nil != err {
			return n, err
		}

		n += i
	}

	return n, nil
}

// WriteLengthPrefixedMessage writes a message to the supplied writer. It prefixes the message with a four byte little-endian representation of the size of the message. The total number of bytes written is returned, along with any error arising.
func WriteLengthPrefixedMessage(w io.Writer, pb proto.Message) (int, error) {
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

	return LengthPrefixFramedMessageSize(pb), nil
}
