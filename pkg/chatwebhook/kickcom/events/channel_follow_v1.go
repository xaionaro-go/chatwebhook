package events

import (
	"github.com/xaionaro-go/chatwebhook/pkg/chatwebhook/kickcom/structs"
	"github.com/xaionaro-go/chatwebhook/pkg/grpc/protobuf/go/chatwebhook_grpc"
)

type ChannelFollowedV1 struct {
	Broadcaster structs.UserV1 `json:"broadcaster"`
	Follower    structs.UserV1 `json:"follower"`
}

func (ChannelFollowedV1) Version() int {
	return 1
}

func (ChannelFollowedV1) TypeName() string {
	return "channel.followed"
}

func (ChannelFollowedV1) TypeID() chatwebhook_grpc.PlatformEventType {
	return chatwebhook_grpc.PlatformEventType_platformEventTypeFollow
}

func (ev *ChannelFollowedV1) ToGRPC() []*chatwebhook_grpc.Event {
	return []*chatwebhook_grpc.Event{{
		Id:            randomEventID(),
		EventType:     ev.TypeID(),
		User:          userToGRPC(ev.Follower),
		TargetChannel: userToGRPC(ev.Broadcaster),
	}}
}
