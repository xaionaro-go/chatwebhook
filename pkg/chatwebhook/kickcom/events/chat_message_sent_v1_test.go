package events

import (
	"reflect"
	"testing"
	"time"

	"github.com/xaionaro-go/chatwebhook/pkg/chatwebhook/kickcom/structs"
	"github.com/xaionaro-go/chatwebhook/pkg/grpc/protobuf/go/chatwebhook_grpc"
)

func TestChatMessageSentV1_ToGRPC(t *testing.T) {
	randomEventID = func() string {
		return "random-event-id"
	}
	tests := []struct {
		name string
		ev   ChatMessageSentV1
		want []*chatwebhook_grpc.Event
	}{
		{
			name: "simple message without emotes",
			ev: ChatMessageSentV1{
				MessageID: "msg1",
				Content:   "Hello.. [emote]; Hello!",
				Sender: structs.UserV1{
					UserID:   1,
					Username: "testuser",
				},
				Broadcaster: structs.UserV1{
					UserID:   2,
					Username: "broadcaster",
				},
				CreatedAt: "2025-11-02T00:00:00Z",
				Emotes: []structs.EmoteV1{{
					EmoteID:   "emote-id",
					Positions: []structs.PositionV1{{S: 8, E: 15}},
				}},
			},
			want: []*chatwebhook_grpc.Event{
				{
					Id:                "msg1",
					CreatedAtUNIXNano: uint64(time.Date(2025, time.November, 2, 0, 0, 0, 0, time.UTC).UnixNano()),
					EventType:         chatwebhook_grpc.PlatformEventType_platformEventTypeChatMessage,
					User: &chatwebhook_grpc.User{
						Id:   "1",
						Name: "testuser",
					},
					TargetChannel: &chatwebhook_grpc.User{
						Id:   "2",
						Name: "broadcaster",
					},
					Message: &chatwebhook_grpc.Message{
						Content:    "Hello.. <img class=\"kick-emote\" src=\"https://files.kick.com/emotes/emote-id/fullsize\" alt=\":emote:\">; Hello!",
						FormatType: chatwebhook_grpc.TextFormatType_TEXT_FORMAT_TYPE_HTML,
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.ev.ToGRPC(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ChatMessageSentV1.ToGRPC() = %v, want %v", got, tt.want)
			}
		})
	}
}
