package stream

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/golang/protobuf/proto"
	"github.com/robinpowered/go-proto/collection"
	"github.com/robinpowered/go-proto/message"
)

// LengthPrefixFramedCollectionSize calculates the size of the expected stream produced when encoding this collection using length prefixed framing.
// The implementation assumes a 32 bit varint prefix.
func LengthPrefixFramedCollectionSize(mc collection.MessageCollection) (n int) {
	for _, m := range mc {
		n += LengthPrefixFramedMessageSize(m)
	}

	return n
}

// LengthPrefixFramedMessageSize calculates the expected total size of a length-prefixed framed protobuf message within a stream.
// It performs a simple addition of the length of the prefix field and the protobuf message size.
// The implementation assumes a 32 bit varint prefix.
func LengthPrefixFramedMessageSize(m proto.Message) int {
	_, s := prefix(m)

	return s + proto.Size(m)
}

// prefix returns the prefix frame of a proto message, plus its encoded length.
func prefix(m proto.Message) ([binary.MaxVarintLen32]byte, int) {
	var buf [binary.MaxVarintLen32]byte
	encodedLength := binary.PutUvarint(buf[:], uint64(proto.Size(m)))

	return buf, encodedLength
}

// ReadLengthPrefixedCollection reads a collection of protocol buffer messages from the supplied reader.
// Each message is presumed prefixed by a 32 bit varint which represents the size of the ensuing message.
// The UnmarshalFunc argument is a supplied callback used to convert the raw bytes read as a message to the desired message type.
// The protocol buffer message collection is returned, along with any error arising.
// For more detailed information on this approach, see the official protocol buffer documentation https://developers.google.com/protocol-buffers/docs/techniques#streaming.
func ReadLengthPrefixedCollection(r io.Reader, f message.UnmarshalFunc) (pbs collection.MessageCollection, err error) {
	for {
		var prefixBuf [binary.MaxVarintLen32]byte
		var bytesRead, varIntBytes int
		var messageLength uint64
		for varIntBytes == 0 { // i.e. no varint has been decoded yet.
			if bytesRead >= len(prefixBuf) {
				return pbs, fmt.Errorf("invalid varint32 encountered")
			}
			// We have to read byte by byte here to avoid reading more bytes
			// than required. Each read byte is appended to what we have
			// read before.
			newBytesRead, err := r.Read(prefixBuf[bytesRead : bytesRead+1])
			if newBytesRead == 0 {
				if io.EOF == err {
					return pbs, nil
				} else if err != nil {
					return pbs, err
				}
				// A Reader should not return (0, nil), but if it does,
				// it should be treated as no-op (according to the
				// Reader contract). So let's go on...
				continue
			}
			bytesRead += newBytesRead
			// Now present everything read so far to the varint decoder and
			// see if a varint can be decoded already.
			messageLength, varIntBytes = proto.DecodeVarint(prefixBuf[:bytesRead])
		}

		messageBuf := make([]byte, messageLength)
		newBytesRead, err := io.ReadFull(r, messageBuf)
		bytesRead += newBytesRead
		if err != nil {
			return pbs, err
		}

		pb, err := f(messageBuf)
		if nil != err {
			return nil, err
		}

		pbs = append(pbs, pb)
	}
}

// WriteLengthPrefixedCollection writes the collection of protobuf messages to the supplied writer.
// The write operation uses length-prefixed framing. This means that each protocol buffer message is prefixed
// by its length. This implementation encodes the length as a 32 bit varint prefix.
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

// WriteLengthPrefixedMessage writes a message to the supplied writer. It prefixes the message with a 32 bit varint representation of the size of the message. The total number of bytes written is returned, along with any error arising.
func WriteLengthPrefixedMessage(w io.Writer, pb proto.Message) (int, error) {
	b, err := proto.Marshal(pb)
	if nil != err {
		return 0, err
	}

	prefix, s := prefix(pb)

	sync, err := w.Write(prefix[:s])
	if nil != err {
		return sync, err
	}

	n, err := w.Write(b)

	return n + sync, err
}
