package events

import (
	"github.com/xaionaro-go/chatwebhook/pkg/chatwebhook/kickcom/structs"
	"github.com/xaionaro-go/chatwebhook/pkg/grpc/protobuf/go/chatwebhook_grpc"
)

type ChannelSubscriptionGiftsV1 struct {
	Broadcaster structs.UserV1   `json:"broadcaster"`
	Gifter      structs.UserV1   `json:"gifter"`
	Giftees     []structs.UserV1 `json:"giftees"`
	CreatedAt   string           `json:"created_at"`
	ExpiresAt   string           `json:"expires_at"`
}

func (ChannelSubscriptionGiftsV1) Version() int {
	return 1
}

func (ChannelSubscriptionGiftsV1) TypeName() string {
	return "channel.subscription.gifts"
}

func (ChannelSubscriptionGiftsV1) TypeID() chatwebhook_grpc.PlatformEventType {
	return chatwebhook_grpc.PlatformEventType_platformEventTypeGiftedSubscription
}

func (e *ChannelSubscriptionGiftsV1) ToGRPC() []*chatwebhook_grpc.Event {
	result := make([]*chatwebhook_grpc.Event, 0, len(e.Giftees))
	for _, giftee := range e.Giftees {
		result = append(result, &chatwebhook_grpc.Event{
			Id:                randomEventID(),
			CreatedAtUNIXNano: timeToGRPC(e.CreatedAt),
			ExpiresAtUNIXNano: ptr(timeToGRPC(e.ExpiresAt)),
			EventType:         e.TypeID(),
			User:              userToGRPC(e.Gifter),
			TargetUser:        userToGRPC(giftee),
			TargetChannel:     userToGRPC(e.Broadcaster),
		})
	}
	return result
}
