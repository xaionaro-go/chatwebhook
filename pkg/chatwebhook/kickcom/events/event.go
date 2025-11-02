package events

import (
	"encoding/json"
	"fmt"

	"github.com/xaionaro-go/chatwebhook/pkg/grpc/protobuf/go/chatwebhook_grpc"
)

const (
	HTTPHeaderEventType    = "Kick-Event-Type"
	HTTPHeaderEventVersion = "Kick-Event-Version"
)

type Event interface {
	Version() int
	TypeName() string
	TypeID() chatwebhook_grpc.PlatformEventType
	ToGRPC() []*chatwebhook_grpc.Event
}

func All() []Event {
	return []Event{
		&ChannelFollowedV1{},
		&ChannelSubscriptionGiftsV1{},
		&ChannelSubscriptionNewV1{},
		&ChannelSubscriptionRenewalV1{},
		&ChatMessageSentV1{},
		&KicksGifted{},
		&LiveStreamMetadataUpdatedV1{},
		&LiveStreamStatusUpdatedV1{},
		&ModerationBannedV1{},
	}
}

func AbstractParse(eventType string, version int, eventJSON []byte) (Event, error) {
	if len(eventJSON) == 0 {
		return nil, fmt.Errorf("empty event JSON for event type %q version %d", eventType, version)
	}
	for _, ev := range All() {
		if ev.TypeName() != eventType || ev.Version() != version {
			continue
		}
		err := json.Unmarshal(eventJSON, ev)
		if err != nil {
			return nil, fmt.Errorf("unable to unmarshal event %q of type %q and version %d: %w", eventJSON, eventType, version, err)
		}
		return ev, nil
	}
	return nil, fmt.Errorf("unknown event type %q version %d", eventType, version)
}

func Parse[E Event](data []byte, evn E) (E, error) {
	err := json.Unmarshal(data, &evn)
	if err != nil {
		return evn, err
	}
	return evn, nil
}
