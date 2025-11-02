package events

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/xaionaro-go/chatwebhook/pkg/chatwebhook/kickcom/structs"
	"github.com/xaionaro-go/chatwebhook/pkg/grpc/protobuf/go/chatwebhook_grpc"
)

var randomEventID = func() string {
	return must(uuid.NewRandom()).String()
}

func userToGRPC(u structs.UserV1) *chatwebhook_grpc.User {
	return &chatwebhook_grpc.User{
		Id:   fmt.Sprintf("%d", u.UserID),
		Slug: u.ChannelSlug,
		Name: u.Username,
	}
}

func timeToGRPC(timeStr structs.RFC3339String) uint64 {
	return uint64(must(timeStr.Parse()).UnixNano())
}
