package pbs

import (
	"bytes"
	"testing"

	"github.com/golang/protobuf/proto"
	"github.com/mcos/protobuf-collections/mock"
	"github.com/stretchr/testify/assert"
)

var (
	mockMessages = MessageCollection{
		&mock.Message{Name: "Foo"},
		&mock.Message{Name: "Bar"},
		&mock.Message{Name: "Baz"},
	}

	mockUnmarshalFunc UnmarshalFunc = func(b []byte) (proto.Message, error) {
		var m mock.Message
		err := proto.Unmarshal(b, &m)

		return &m, err
	}
)

func TestLengthPrefixFramedSize(t *testing.T) {
	s := 0
	for _, m := range mockMessages {
		s += LengthPrefixedFramedSize(m)
	}

	assert.Equal(t, s, mockMessages.LengthPrefixFramedSize(), "they should be equal")
}

func TestWriteLengthPrefixedMessage(t *testing.T) {
	buf := new(bytes.Buffer)

	pb := &mock.Message{Name: "Single Message"}

	n, err := writeLengthPrefixedMessage(buf, pb)

	assert.Equal(t, n, LengthPrefixedFramedSize(pb), "they should be equal")
	assert.Nil(t, err)
}

func TestWriteLengthPrefixedCollection(t *testing.T) {
	buf := new(bytes.Buffer)

	n, err := WriteLengthPrefixedCollection(buf, mockMessages)

	assert.Equal(t, n, buf.Len(), "they should be equal")
	assert.Equal(t, n, mockMessages.LengthPrefixFramedSize(), "they should be equal")
	assert.Nil(t, err)
}

func TestWriteLengthPrefixedCollectionFails(t *testing.T) {
	buf := new(bytes.Buffer)

	_, _ = buf.Write([]byte("This will fail because this string messes up the encoding"))

	n, err := WriteLengthPrefixedCollection(buf, mockMessages)

	assert.NotEqual(t, n, buf.Len(), "they should not be equal")
	assert.NotEqual(t, buf.Len(), mockMessages.LengthPrefixFramedSize(), "they should not be equal")
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

	pbs, err := ReadLengthPrefixedCollection(buf, mockUnmarshalFunc)

	assert.NotEqual(t, len(pbs), len(mockMessages))
	assert.NotNil(t, err)
}
