package client

import (
	"context"
	"fmt"
	"io"

	"github.com/facebookincubator/go-belt/tool/logger"
	"github.com/xaionaro-go/chatwebhook/pkg/grpc/protobuf/go/chatwebhook_grpc"
	"github.com/xaionaro-go/chatwebhook/pkg/xgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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
) (chatwebhook_grpc.ChatWebHookClient, io.Closer, error) {
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
	logger.Errorf(ctx, "gRPC call error: %v", err)
	return err
}

func (c *Client) GetCallWrapper() xgrpc.CallWrapperFunc {
	return nil
}

func (c *Client) GetMessagesChan(
	ctx context.Context,
	platformID chatwebhook_grpc.PlatformID,
	channelID string,
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
					ChannelID:  channelID,
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
