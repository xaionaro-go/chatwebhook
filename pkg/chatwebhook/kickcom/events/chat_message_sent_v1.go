package events

import (
	"fmt"
	"html"
	"regexp"
	"slices"
	"strings"

	"github.com/xaionaro-go/chatwebhook/pkg/chatwebhook/kickcom/structs"
	"github.com/xaionaro-go/chatwebhook/pkg/grpc/protobuf/go/chatwebhook_grpc"
)

const (
	forceParseEmotes = true
)

type ChatMessageSentV1 struct {
	MessageID   string               `json:"message_id"`
	RepliesTo   *structs.RepliesToV1 `json:"replies_to"`
	Broadcaster structs.UserV1       `json:"broadcaster"`
	Sender      structs.UserV1       `json:"sender"`
	Content     string               `json:"content"`
	Emotes      []structs.EmoteV1    `json:"emotes"`
	CreatedAt   string               `json:"created_at"`
}

func (ChatMessageSentV1) Version() int {
	return 1
}

func (ChatMessageSentV1) TypeName() string {
	return "chat.message.sent"
}

func (ChatMessageSentV1) TypeID() chatwebhook_grpc.PlatformEventType {
	return chatwebhook_grpc.PlatformEventType_platformEventTypeChatMessage
}

var emoteRegexp = regexp.MustCompile(`\[emote:([0-9]+):[^\]]+\]`)

func (ev *ChatMessageSentV1) ToGRPC() []*chatwebhook_grpc.Event {
	type emoteIDWithPosition struct {
		EmoteID  string
		Position structs.PositionV1
	}
	var emotesWithPositions []emoteIDWithPosition
	for _, emote := range ev.Emotes {
		for _, position := range emote.Positions {
			emotesWithPositions = append(emotesWithPositions, emoteIDWithPosition{
				EmoteID:  emote.EmoteID,
				Position: position,
			})
		}
	}
	slices.SortFunc(emotesWithPositions, func(a, b emoteIDWithPosition) int {
		switch {
		case a.Position.S < b.Position.S:
			return -1
		case a.Position.S > b.Position.S:
			return 1
		default:
			return 0
		}
	})
	if forceParseEmotes && len(emotesWithPositions) == 0 {
		matches := emoteRegexp.FindAllStringSubmatchIndex(ev.Content, -1)
		for _, m := range matches {
			emoteID := ev.Content[m[2]:m[3]]
			emotesWithPositions = append(emotesWithPositions, emoteIDWithPosition{
				EmoteID:  emoteID,
				Position: structs.PositionV1{S: int64(m[0]), E: int64(m[1])},
			})
		}
	}

	var result strings.Builder
	inIdx := 0
	for _, emote := range emotesWithPositions {
		result.WriteString(html.EscapeString(ev.Content[inIdx:emote.Position.S]))
		result.WriteString(fmt.Sprintf(`<img class="kick-emote" src="https://files.kick.com/emotes/%s/fullsize" alt=":emote:">`, emote.EmoteID))
		inIdx = int(emote.Position.E)
	}
	result.WriteString(html.EscapeString(ev.Content[inIdx:len(ev.Content)]))
	var inReplyTo *string
	if ev.RepliesTo != nil {
		inReplyTo = ptr(ev.RepliesTo.MessageID)
	}
	return []*chatwebhook_grpc.Event{
		{
			Id:                ev.MessageID,
			CreatedAtUNIXNano: timeToGRPC(ev.CreatedAt),
			EventType:         ev.TypeID(),
			User:              userToGRPC(ev.Sender),
			TargetChannel:     userToGRPC(ev.Broadcaster),
			Message: &chatwebhook_grpc.Message{
				Content:    result.String(),
				FormatType: chatwebhook_grpc.TextFormatType_TEXT_FORMAT_TYPE_HTML,
				InReplyTo:  inReplyTo,
			},
		},
	}
}
