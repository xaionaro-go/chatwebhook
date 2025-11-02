package events

import (
	"github.com/xaionaro-go/chatwebhook/pkg/chatwebhook/kickcom/structs"
	"github.com/xaionaro-go/chatwebhook/pkg/grpc/protobuf/go/chatwebhook_grpc"
)

type ChannelSubscriptionNewV1 struct {
	Broadcaster structs.UserV1 `json:"broadcaster"`
	Subscriber  structs.UserV1 `json:"subscriber"`
	Duration    int64          `json:"duration"`
	CreatedAt   string         `json:"created_at"`
	ExpiresAt   string         `json:"expires_at"`
}

func (ChannelSubscriptionNewV1) Version() int {
	return 1
}

func (ChannelSubscriptionNewV1) TypeName() string {
	return "channel.subscription.new"
}

func (ChannelSubscriptionNewV1) TypeID() chatwebhook_grpc.PlatformEventType {
	return chatwebhook_grpc.PlatformEventType_platformEventTypeSubscriptionNew
}

func (e ChannelSubscriptionNewV1) ToGRPC() []*chatwebhook_grpc.Event {
	return []*chatwebhook_grpc.Event{
		{
			Id:                randomEventID(),
			CreatedAtUNIXNano: timeToGRPC(e.CreatedAt),
			ExpiresAtUNIXNano: ptr(timeToGRPC(e.ExpiresAt)),
			EventType:         e.TypeID(),
			User:              userToGRPC(e.Subscriber),
			TargetChannel:     userToGRPC(e.Broadcaster),
		},
	}
}
