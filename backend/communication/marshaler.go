package communication

import (
	"github.com/golang/protobuf/proto"
)

func marshalMessage(message proto.Message) *[]byte {
	bytes, err := proto.Marshal(message)
	if err != nil {
		panic(err)
	}

	return &bytes
}
