package events

import (
	"encoding/json"

	"github.com/xaionaro-go/chatwebhook/pkg/chatwebhook/kickcom/structs"
	"github.com/xaionaro-go/chatwebhook/pkg/grpc/protobuf/go/chatwebhook_grpc"
)

type KicksGifted struct {
	Broadcaster structs.UserV1 `json:"broadcaster"`
	Sender      structs.UserV1 `json:"sender"`
	Gift        structs.GiftV1 `json:"gift"`
	CreatedAt   string         `json:"created_at"`
}

func (KicksGifted) Version() int {
	return 1
}

func (KicksGifted) TypeName() string {
	return "kicks.gifted"
}

func (KicksGifted) TypeID() chatwebhook_grpc.PlatformEventType {
	return chatwebhook_grpc.PlatformEventType_platformEventTypeOther
}

func (ev *KicksGifted) ToGRPC() []*chatwebhook_grpc.Event {
	return []*chatwebhook_grpc.Event{{
		Id:                randomEventID(),
		CreatedAtUNIXNano: timeToGRPC(ev.CreatedAt),
		EventType:         ev.TypeID(),
		User:              userToGRPC(ev.Sender),
		TargetChannel:     userToGRPC(ev.Broadcaster),
		Tier:              ptr(ev.Gift.Tier),
		Message: &chatwebhook_grpc.Message{
			Content:    string(must(json.Marshal(ev.Gift))),
			FormatType: chatwebhook_grpc.TextFormatType_TEXT_FORMAT_TYPE_PLAIN,
		},
	}}
}
