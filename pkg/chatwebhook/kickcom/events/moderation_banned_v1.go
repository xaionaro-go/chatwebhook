package events

import (
	"github.com/xaionaro-go/chatwebhook/pkg/chatwebhook/kickcom/structs"
	"github.com/xaionaro-go/chatwebhook/pkg/grpc/protobuf/go/chatwebhook_grpc"
)

type ModerationBannedV1 struct {
	Broadcaster structs.UserV1        `json:"broadcaster"`
	Moderator   structs.UserV1        `json:"moderator"`
	BannedUser  structs.UserV1        `json:"banned_user"`
	Metadata    structs.BanMetadataV1 `json:"metadata"`
}

func (ModerationBannedV1) Version() int {
	return 1
}

func (ModerationBannedV1) TypeName() string {
	return "moderation.banned"
}

func (ModerationBannedV1) TypeID() chatwebhook_grpc.PlatformEventType {
	return chatwebhook_grpc.PlatformEventType_platformEventTypeBan
}

func (ev *ModerationBannedV1) ToGRPC() []*chatwebhook_grpc.Event {
	return []*chatwebhook_grpc.Event{{
		Id:            randomEventID(),
		EventType:     ev.TypeID(),
		User:          userToGRPC(ev.Moderator),
		TargetChannel: userToGRPC(ev.Broadcaster),
		TargetUser:    userToGRPC(ev.BannedUser),
	}}
}
