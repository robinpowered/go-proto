package stream

import (
	"bytes"
	"testing"

	"github.com/robinpowered/go-proto/collection"
	"github.com/robinpowered/go-proto/message"
	"github.com/robinpowered/go-proto/mock"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/proto"
)

var (
	mockMessages = collection.MessageCollection{
		&mock.Message{Name: "Foo"},
		&mock.Message{Name: "Bar"},
		&mock.Message{Name: "Baz"},
	}

	mockUnmarshalFunc message.UnmarshalFunc = func(b []byte) (proto.Message, error) {
		var m mock.Message
		err := proto.Unmarshal(b, &m)

		return &m, err
	}
)

func TestLengthPrefixFramedSize(t *testing.T) {
	s := 0
	for _, m := range mockMessages {
		s += LengthPrefixFramedMessageSize(m)
	}

	assert.Equal(t, s, LengthPrefixFramedCollectionSize(mockMessages), "they should be equal")
}

func TestWriteLengthPrefixedMessage(t *testing.T) {
	buf := new(bytes.Buffer)

	pb := &mock.Message{Name: "Single Message"}

	n, err := WriteLengthPrefixedMessage(buf, pb)

	assert.Equal(t, n, LengthPrefixFramedMessageSize(pb), "they should be equal")
	assert.Nil(t, err)
}

func TestWriteLengthPrefixedCollection(t *testing.T) {
	buf := new(bytes.Buffer)

	n, err := WriteLengthPrefixedCollection(buf, mockMessages)

	assert.Equal(t, n, buf.Len(), "they should be equal")
	assert.Equal(t, n, LengthPrefixFramedCollectionSize(mockMessages), "they should be equal")
	assert.Nil(t, err)
}

func TestWriteLengthPrefixedCollectionFails(t *testing.T) {
	buf := new(bytes.Buffer)

	_, _ = buf.Write([]byte("This will fail because this string messes up the encoding"))

	n, err := WriteLengthPrefixedCollection(buf, mockMessages)

	assert.NotEqual(t, n, buf.Len(), "they should not be equal")
	assert.NotEqual(t, buf.Len(), LengthPrefixFramedCollectionSize(mockMessages), "they should not be equal")
	assert.Nil(t, err)
}

func TestReadLengthPrefixedCollection(t *testing.T) {
	buf := new(bytes.Buffer)

	_, err := WriteLengthPrefixedCollection(buf, mockMessages)

	assert.Nil(t, err)

	pbs, err := ReadLengthPrefixedCollection(buf, mockUnmarshalFunc)

	assert.Equal(t, len(pbs), len(mockMessages))
	assert.Nil(t, err)

	for i, pb := range pbs {
		assert.True(t, proto.Equal(pb, mockMessages[i]), "they should be equal")
	}
}

func TestReadLengthPrefixedCollectionFails(t *testing.T) {
	buf := new(bytes.Buffer)

	_, err := WriteLengthPrefixedCollection(buf, mockMessages)

	assert.Nil(t, err)

	_, _ = buf.Write([]byte("This string messes up the decoding"))

	_, err = ReadLengthPrefixedCollection(buf, mockUnmarshalFunc)

	assert.NotNil(t, err)
}
