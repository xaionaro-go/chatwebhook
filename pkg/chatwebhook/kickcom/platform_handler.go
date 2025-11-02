package kickcom

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/xaionaro-go/chatwebhook/pkg/cache"
	"github.com/xaionaro-go/chatwebhook/pkg/chatwebhook"
	"github.com/xaionaro-go/chatwebhook/pkg/chatwebhook/kickcom/events"
	"github.com/xaionaro-go/chatwebhook/pkg/grpc/protobuf/go/chatwebhook_grpc"
)

type PlatformHandler struct {
	PubKeyVerifier *PubKeyVerifier
}

var _ chatwebhook.PlatformHandler = (*PlatformHandler)(nil)

func NewPlatformHandler(
	ctx context.Context,
	cache cache.Cache,
) (*PlatformHandler, error) {
	pubKeyVerifier, err := NewPubKeyVerifier(ctx, cache)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize kick.com public key verifier: %w", err)
	}
	return &PlatformHandler{
		PubKeyVerifier: pubKeyVerifier,
	}, nil
}

func (h *PlatformHandler) PlatformID() chatwebhook_grpc.PlatformID {
	return ID
}

func (h *PlatformHandler) ParseEvents(
	r *http.Request,
) ([]*chatwebhook_grpc.Event, error) {
	if err := h.PubKeyVerifier.VerifyRequest(r); err != nil {
		return nil, fmt.Errorf("unable to verify kick.com request signature: %w", err)
	}

	return parseEvents(r)
}

func parseEvents(
	r *http.Request,
) ([]*chatwebhook_grpc.Event, error) {
	eventType := r.Header.Get(events.HTTPHeaderEventType)
	if eventType == "" {
		return nil, fmt.Errorf("missing %q header", events.HTTPHeaderEventType)
	}

	eventVersionStr := r.Header.Get(events.HTTPHeaderEventVersion)
	if eventVersionStr == "" {
		return nil, fmt.Errorf("missing %q header", events.HTTPHeaderEventVersion)
	}
	eventVersion, err := strconv.ParseUint(eventVersionStr, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("unable to parse %q header %q: %w", events.HTTPHeaderEventVersion, eventVersionStr, err)
	}

	data, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, fmt.Errorf("unable to read request body: %w", err)
	}

	ev, err := events.AbstractParse(eventType, int(eventVersion), data)
	if err != nil {
		return nil, fmt.Errorf("unable to parse kick.com event: %w", err)
	}

	return ev.ToGRPC(), nil
}
