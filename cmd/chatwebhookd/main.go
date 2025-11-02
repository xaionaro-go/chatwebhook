package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"strconv"

	"github.com/facebookincubator/go-belt"
	"github.com/facebookincubator/go-belt/tool/logger"
	"github.com/facebookincubator/go-belt/tool/logger/implementation/zap"
	"github.com/xaionaro-go/chatwebhook/pkg/cache/cachedir"
	"github.com/xaionaro-go/chatwebhook/pkg/chatwebhook"
	"github.com/xaionaro-go/chatwebhook/pkg/chatwebhook/kickcom"
	"github.com/xaionaro-go/chatwebhook/pkg/grpc/server"
	"github.com/xaionaro-go/observability"
	"github.com/xaionaro-go/xpath"
)

func newListener(addr string) (net.Listener, error) {
	host, portString, err := net.SplitHostPort(addr)
	if err != nil {
		return nil, fmt.Errorf("unable to parse address %s: %w", addr, err)
	}

	ip := net.ParseIP(host)
	port, err := strconv.ParseUint(portString, 10, 16)
	if err != nil {
		return nil, fmt.Errorf("unable to parse port %s: %w", portString, err)
	}

	listener, err := net.ListenTCP("tcp", &net.TCPAddr{
		IP:   ip,
		Port: int(port),
		Zone: "",
	})
	if err != nil {
		return nil, fmt.Errorf("unable to listen on %s: %w", addr, err)
	}

	return listener, nil
}

var void = struct{}{}

func main() {
	loggerLevel := logger.LevelInfo
	flag.Var(&loggerLevel, "log-level", "")
	subscribeAddr := flag.String("subscribe-addr", ":8081", "Address of the subscribing to the messages")
	receiverAddr := flag.String("receiver-addr", ":8080", "Address to listen on for the messages")
	noTLSFlag := flag.Bool("no-tls", false, "Disable TLS for the receiver")
	certPath := flag.String("cert", "/etc/chatwebhook/server.crt", "Path to TLS certificate")
	keyPath := flag.String("key", "/etc/chatwebhook/server.key", "Path to TLS key")
	cacheDirFlag := flag.String("cache-dir", "~/.local/cache/chatwebhook", "Path to cache directory")
	flag.Parse()
	l := zap.Default().WithLevel(loggerLevel)
	ctx := context.Background()
	ctx = logger.CtxWithLogger(ctx, l)
	logger.Default = func() logger.Logger {
		return l
	}
	defer belt.Flush(ctx)

	cacheDir := must(xpath.Expand(*cacheDirFlag))
	logger.Debugf(ctx, "using cache dir: %s", cacheDir)
	must(void, os.MkdirAll(cacheDir, 0o755))

	subscribeListener := must(newListener(*subscribeAddr))
	receiverListener := must(newListener(*receiverAddr))
	cache := must(cachedir.New(cacheDir))

	h := chatwebhook.NewHandler()
	h.SetPlatformHandler(must(kickcom.NewPlatformHandler(ctx, cache)))

	errCh := make(chan error, 1)

	observability.Go(ctx, func(ctx context.Context) {
		grpcSrv := server.New(h)
		logger.Infof(ctx, "starting subscription gRPC server at %s", subscribeListener.Addr())
		err := grpcSrv.Serve(subscribeListener)
		errCh <- fmt.Errorf("subscription gRPC server error: %w", err)
	})

	observability.Go(ctx, func(ctx context.Context) {
		mux := http.NewServeMux()
		mux.HandleFunc("/kick", h.GetPublishFunc(kickcom.ID))
		srv := &http.Server{
			Handler: mux,
		}
		if *noTLSFlag {
			logger.Infof(ctx, "starting receiver HTTP server at %s without TLS", receiverListener.Addr())
			err := srv.Serve(receiverListener)
			errCh <- fmt.Errorf("receiver HTTP server error: %w", err)
			return
		}
		logger.Infof(ctx, "starting receiver HTTP server at %s with TLS", receiverListener.Addr())
		err := srv.ServeTLS(receiverListener, *certPath, *keyPath)
		errCh <- fmt.Errorf("receiver HTTP server error: %w", err)
	})

	panic(<-errCh)
}
