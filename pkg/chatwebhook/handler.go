package chatwebhook

import (
	"bytes"
	"context"
	"io"
	"net/http"

	"github.com/facebookincubator/go-belt/tool/logger"
	"github.com/xaionaro-go/chatwebhook/pkg/grpc/protobuf/go/chatwebhook_grpc"
	"github.com/xaionaro-go/eventbus"
	"github.com/xaionaro-go/secret"
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
		if r.GetBody == nil {
			body, err := io.ReadAll(r.Body)
			if err != nil {
				http.Error(w, "unable to read request body: "+err.Error(), http.StatusBadRequest)
				return
			}
			r.GetBody = func() (io.ReadCloser, error) {
				return io.NopCloser(bytes.NewReader(body)), nil
			}
		}

		bodyReader, err := r.GetBody()
		if err != nil {
			http.Error(w, "unable to get request body: "+err.Error(), http.StatusBadRequest)
			return
		}
		defer bodyReader.Close()

		body, err := io.ReadAll(bodyReader)
		if err != nil {
			http.Error(w, "unable to read request body: "+err.Error(), http.StatusBadRequest)
			return
		}

		logger.Tracef(ctx, "handling webhook request for platform %s (body: %s)", platformID.String(), string(body))

		apiKey := r.Form.Get("apiKey")
		platformHandler := h.GetPlatformHandler(platformID)
		if platformHandler == nil {
			http.Error(w, "unsupported platform", http.StatusBadRequest)
			return
		}

		logger.Debugf(ctx, "received webhook request for platform %s", platformID.String())

		events, err := platformHandler.ParseEvents(r)
		if err != nil {
			http.Error(w, "unable to parse event: "+err.Error(), http.StatusBadRequest)
			return
		}

		for _, event := range events {
			subKey := subKey{
				PlatformID: platformID,
				APIKey:     apiKey,
			}
			logger.Debugf(ctx, "publishing event %+v to subscribers of %s", event, platformID.String())
			eventbus.SendEventWithCustomTopic(ctx, h.EventBus, subKey, event)
		}
		w.WriteHeader(http.StatusOK)
	}
}

type subKey struct {
	PlatformID chatwebhook_grpc.PlatformID
	APIKey     string
}

func (h *Handler) Subscribe(
	ctx context.Context,
	platformID chatwebhook_grpc.PlatformID,
	apiKey secret.String,
) (<-chan *chatwebhook_grpc.Event, error) {
	return eventbus.SubscribeWithCustomTopic[subKey, *chatwebhook_grpc.Event](ctx,
		h.EventBus, subKey{
			PlatformID: platformID,
			APIKey:     apiKey.Get(),
		},
	).EventChan(), nil
}
