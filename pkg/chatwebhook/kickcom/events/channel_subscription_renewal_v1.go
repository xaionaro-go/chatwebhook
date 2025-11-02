package events

import (
	"github.com/google/uuid"
	"github.com/xaionaro-go/chatwebhook/pkg/chatwebhook/kickcom/structs"
	"github.com/xaionaro-go/chatwebhook/pkg/grpc/protobuf/go/chatwebhook_grpc"
)

type ChannelSubscriptionRenewalV1 struct {
	Broadcaster structs.UserV1        `json:"broadcaster"`
	Subscriber  structs.UserV1        `json:"subscriber"`
	Duration    int64                 `json:"duration"`
	CreatedAt   structs.RFC3339String `json:"created_at"`
	ExpiresAt   structs.RFC3339String `json:"expires_at"`
}

func (ChannelSubscriptionRenewalV1) Version() int {
	return 1
}

func (ChannelSubscriptionRenewalV1) TypeName() string {
	return "channel.subscription.renewal"
}

func (ChannelSubscriptionRenewalV1) TypeID() chatwebhook_grpc.PlatformEventType {
	return chatwebhook_grpc.PlatformEventType_platformEventTypeSubscriptionRenewed
}

func (e *ChannelSubscriptionRenewalV1) ToGRPC() []*chatwebhook_grpc.Event {
	return []*chatwebhook_grpc.Event{
		{
			Id:                must(uuid.NewRandom()).String(),
			CreatedAtUNIXNano: timeToGRPC(e.CreatedAt),
			ExpiresAtUNIXNano: ptr(timeToGRPC(e.ExpiresAt)),
			EventType:         e.TypeID(),
			User:              userToGRPC(e.Subscriber),
			TargetChannel:     userToGRPC(e.Broadcaster),
		},
	}
}
