package server

import (
	"context"
	"net"

	"github.com/xaionaro-go/chatwebhook/pkg/chatwebhook"
	"github.com/xaionaro-go/chatwebhook/pkg/grpc/protobuf/go/chatwebhook_grpc"
	"github.com/xaionaro-go/chatwebhook/pkg/xgrpc"
	"google.golang.org/grpc"
)

type Server struct {
	chatwebhook_grpc.UnimplementedChatWebHookServer
	Handler    *chatwebhook.Handler
	GRPCServer *grpc.Server
}

func New(
	handler *chatwebhook.Handler,
) *Server {
	srv := &Server{
		GRPCServer: grpc.NewServer(),
		Handler:    handler,
	}
	chatwebhook_grpc.RegisterChatWebHookServer(srv.GRPCServer, srv)
	return srv
}

func parseEvent(in *chatwebhook_grpc.Event) *chatwebhook_grpc.Event {
	return in
}

func (srv *Server) Serve(
	listener net.Listener,
) error {
	return srv.GRPCServer.Serve(listener)
}

func (srv *Server) Subscribe(
	req *chatwebhook_grpc.SubscribeRequest,
	c chatwebhook_grpc.ChatWebHook_SubscribeServer,
) error {
	return xgrpc.WrapChan(c.Context(),
		func(ctx context.Context) (<-chan *chatwebhook_grpc.Event, error) {
			return srv.Handler.Subscribe(
				ctx,
				req.PlatformID,
				req.ChannelID,
				req.ApiKey,
			)
		},
		c,
		parseEvent,
	)
}
