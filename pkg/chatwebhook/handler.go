package chatwebhook

import (
	"context"
	"net/http"

	"github.com/xaionaro-go/chatwebhook/pkg/grpc/protobuf/go/chatwebhook_grpc"
	"github.com/xaionaro-go/eventbus"
	"github.com/xaionaro-go/xsync"
)

type PlatformHandler interface {
	PlatformID() chatwebhook_grpc.PlatformID
	ParseEvents(r *http.Request) ([]*chatwebhook_grpc.Event, error)
}

type Handler struct {
	EventBus         *eventbus.EventBus
	PlatformHandlers xsync.Map[chatwebhook_grpc.PlatformID, PlatformHandler]
}

func NewHandler() *Handler {
	return &Handler{
		EventBus: eventbus.New(),
	}
}

func (h *Handler) SetPlatformHandler(
	handler PlatformHandler,
) {
	h.PlatformHandlers.Store(handler.PlatformID(), handler)
}

func (h *Handler) GetPlatformHandler(
	platformID chatwebhook_grpc.PlatformID,
) PlatformHandler {
	v, _ := h.PlatformHandlers.Load(platformID)
	return v
}

func (h *Handler) GetPublishFunc(
	platformID chatwebhook_grpc.PlatformID,
) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		channelID := r.Form.Get("channelID")
		if channelID == "" {
			http.Error(w, "channelID is required", http.StatusBadRequest)
			return
		}

		apiKey := r.Form.Get("apiKey")
		if apiKey == "" {
			http.Error(w, "apiKey is required", http.StatusBadRequest)
			return
		}

		platformHandler := h.GetPlatformHandler(platformID)
		if platformHandler == nil {
			http.Error(w, "unsupported platform", http.StatusBadRequest)
			return
		}

		events, err := platformHandler.ParseEvents(r)
		if err != nil {
			http.Error(w, "unable to parse event: "+err.Error(), http.StatusBadRequest)
			return
		}

		for _, event := range events {
			eventbus.SendEventWithCustomTopic(ctx, h.EventBus, subKey{
				PlatformID: platformID,
				ChannelID:  channelID,
				APIKey:     apiKey,
			}, event)
		}
		w.WriteHeader(http.StatusOK)
	}
}

type subKey struct {
	PlatformID chatwebhook_grpc.PlatformID
	ChannelID  string
	APIKey     string
}

func (h *Handler) Subscribe(
	ctx context.Context,
	platformID chatwebhook_grpc.PlatformID,
	channelID string,
	apiKey string,
) (<-chan *chatwebhook_grpc.Event, error) {
	return eventbus.SubscribeWithCustomTopic[subKey, *chatwebhook_grpc.Event](ctx,
		h.EventBus, subKey{
			PlatformID: platformID,
			ChannelID:  channelID,
			APIKey:     apiKey,
		},
	).EventChan(), nil
}
