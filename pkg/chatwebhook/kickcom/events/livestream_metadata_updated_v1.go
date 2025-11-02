package events

import (
	"encoding/json"

	"github.com/xaionaro-go/chatwebhook/pkg/chatwebhook/kickcom/structs"
	"github.com/xaionaro-go/chatwebhook/pkg/grpc/protobuf/go/chatwebhook_grpc"
)

type LiveStreamMetadataUpdatedV1 struct {
	Broadcaster structs.UserV1     `json:"broadcaster"`
	Metadata    structs.MetadataV1 `json:"metadata"`
}

func (LiveStreamMetadataUpdatedV1) Version() int {
	return 1
}

func (LiveStreamMetadataUpdatedV1) TypeName() string {
	return "livestream.metadata.updated"
}

func (LiveStreamMetadataUpdatedV1) TypeID() chatwebhook_grpc.PlatformEventType {
	return chatwebhook_grpc.PlatformEventType_platformEventTypeStreamInfoUpdate
}

func (ev *LiveStreamMetadataUpdatedV1) ToGRPC() []*chatwebhook_grpc.Event {
	return []*chatwebhook_grpc.Event{{
		Id:            randomEventID(),
		EventType:     ev.TypeID(),
		TargetChannel: userToGRPC(ev.Broadcaster),
		Message: &chatwebhook_grpc.Message{
			Content:    string(must(json.Marshal(ev.Metadata))),
			FormatType: chatwebhook_grpc.TextFormatType_TEXT_FORMAT_TYPE_PLAIN,
		},
	}}
}
