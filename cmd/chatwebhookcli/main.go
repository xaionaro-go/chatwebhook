package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/facebookincubator/go-belt"
	"github.com/facebookincubator/go-belt/tool/logger"
	"github.com/facebookincubator/go-belt/tool/logger/implementation/zap"
	loggertypes "github.com/facebookincubator/go-belt/tool/logger/types"
	"github.com/urfave/cli/v3"
	"github.com/xaionaro-go/chatwebhook/pkg/grpc/client"
)

func getClient(
	ctx context.Context,
	serverAddr string,
) (*client.Client, error) {
	logger.Debugf(ctx, "using server address: %s", serverAddr)
	return client.New(ctx, serverAddr)
}

func main() {
	cmd := &cli.Command{
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "log-level",
				Value: loggertypes.LevelInfo.String(),
				Action: func(ctx context.Context, c *cli.Command, s string) error {
					level, err := loggertypes.ParseLogLevel(s)
					if err != nil {
						return fmt.Errorf("unable to parse log level %q: %w", s, err)
					}
					l := zap.Default().WithLevel(level)
					ctx = logger.CtxWithLogger(ctx, l)
					logger.Default = func() logger.Logger {
						return l
					}
					defer belt.Flush(ctx)

					logger.Debugf(ctx, "using log level: %s", level)
					return c.Action(ctx, c)
				},
			},
			&cli.StringFlag{
				Name:  "server",
				Value: "home.dx.center:8081",
				Usage: "Address of the chatwebhookd gRPC server",
			},
		},
		Commands: []*cli.Command{
			{
				Name: "chat",
				Commands: []*cli.Command{
					{
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:     "platform-id",
								Usage:    "Platform ID to listen to (available values: kick)",
								Required: true,
							},
							&cli.StringFlag{
								Name:     "channel-id",
								Usage:    "Channel ID to listen to",
								Required: true,
							},
							&cli.StringFlag{
								Name:     "api-key",
								Usage:    "API key to access your queue of events",
								Required: true,
							},
						},
						Name: "listen",
						Action: func(ctx context.Context, cmd *cli.Command) error {
							platformID, err := parsePlatformID(cmd.String("platform-id"))
							if err != nil {
								return fmt.Errorf("unable to parse platform ID %q: %w", cmd.String("platform-id"), err)
							}
							logger.Debugf(ctx, "listening to chat messages on %s:%s...", platformID.String(), cmd.String("channel-id"))
							grpcClient := must(getClient(ctx, cmd.String("server")))
							messagesChan, err := grpcClient.GetMessagesChan(
								ctx,
								platformID,
								cmd.String("channel-id"),
								cmd.String("api-key"),
							)
							if err != nil {
								return fmt.Errorf("unable to get messages channel: %w", err)
							}
							for msg := range messagesChan {
								logger.Infof(ctx, "%+v", msg)
							}
							return nil
						},
					},
				},
			},
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
