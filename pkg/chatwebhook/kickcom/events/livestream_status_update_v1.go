package events

import (
	"github.com/xaionaro-go/chatwebhook/pkg/chatwebhook/kickcom/structs"
	"github.com/xaionaro-go/chatwebhook/pkg/grpc/protobuf/go/chatwebhook_grpc"
)

type LiveStreamStatusUpdatedV1 struct {
	Broadcaster structs.UserV1 `json:"broadcaster"`
	IsLive      bool           `json:"is_live"`
	Title       string         `json:"title"`
	StartedAt   string         `json:"started_at"`
	EndedAt     *string        `json:"ended_at"`
}

func (LiveStreamStatusUpdatedV1) Version() int {
	return 1
}

func (LiveStreamStatusUpdatedV1) TypeName() string {
	return "livestream.status.updated"
}

func (ev *LiveStreamStatusUpdatedV1) TypeID() chatwebhook_grpc.PlatformEventType {
	if ev.IsLive {
		return chatwebhook_grpc.PlatformEventType_platformEventTypeStreamOnline
	}
	return chatwebhook_grpc.PlatformEventType_platformEventTypeStreamOffline
}

func (ev LiveStreamStatusUpdatedV1) ToGRPC() []*chatwebhook_grpc.Event {
	return []*chatwebhook_grpc.Event{{
		Id:            randomEventID(),
		EventType:     ev.TypeID(),
		TargetChannel: userToGRPC(ev.Broadcaster),
	}}
}
