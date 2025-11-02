package client

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/facebookincubator/go-belt/tool/logger"
	"github.com/xaionaro-go/chatwebhook/pkg/grpc/protobuf/go/chatwebhook_grpc"
	"github.com/xaionaro-go/xgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

const (
	DefaultServerAddress = "home.dx.center:4531"
)

type Client struct {
	Target string
}

func New(
	ctx context.Context,
	target string,
) (*Client, error) {
	c := &Client{
		Target: target,
	}
	return c, nil
}

func (c *Client) GRPCClient(
	ctx context.Context,
) (_ret0 chatwebhook_grpc.ChatWebHookClient, _ret1 io.Closer, _err error) {
	logger.Tracef(ctx, "GRPCClient")
	defer func() { logger.Tracef(ctx, "/GRPCClient: %v", _err) }()
	conn, err := grpc.NewClient(
		c.Target,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to initialize a gRPC client: %w", err)
	}

	client := chatwebhook_grpc.NewChatWebHookClient(conn)
	return client, conn, nil
}

func (c *Client) ProcessError(ctx context.Context, err error) error {
	logger.Tracef(ctx, "processError(ctx, '%v'): %T", err, err)
	if s, ok := status.FromError(err); ok {
		logger.Tracef(ctx, "processError(ctx, '%v'): code == %#+v; msg == %#+v", err, s.Code(), s.Message())
		switch s.Code() {
		case codes.Unavailable:
			logger.Debugf(ctx, "suppressed the error (forcing a retry in a second)")
			time.Sleep(time.Second)
			return nil
		}
	}
	logger.Errorf(ctx, "gRPC call error: %v", err)
	return err
}

func (c *Client) GetCallWrapper() xgrpc.CallWrapperFunc {
	return nil
}

func (c *Client) GetMessagesChan(
	ctx context.Context,
	platformID chatwebhook_grpc.PlatformID,
	apiKey string,
) (<-chan *chatwebhook_grpc.Event, error) {
	return xgrpc.UnwrapChan(ctx, c,
		func(
			ctx context.Context,
			client chatwebhook_grpc.ChatWebHookClient,
		) (chatwebhook_grpc.ChatWebHook_SubscribeClient, error) {
			return xgrpc.Call(ctx, c,
				client.Subscribe,
				&chatwebhook_grpc.SubscribeRequest{
					PlatformID: platformID,
					ApiKey:     apiKey,
				},
			)
		},
		func(
			ctx context.Context,
			event *chatwebhook_grpc.Event,
		) *chatwebhook_grpc.Event {
			return event
		},
	)
}
