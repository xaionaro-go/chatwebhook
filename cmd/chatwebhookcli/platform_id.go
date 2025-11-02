package main

import (
	"fmt"
	"strings"

	"github.com/xaionaro-go/chatwebhook/pkg/grpc/protobuf/go/chatwebhook_grpc"
)

func parsePlatformID(s string) (platformID chatwebhook_grpc.PlatformID, err error) {
	switch strings.ToLower(s) {
	case "kick", "kickcom", "kick.com":
		return chatwebhook_grpc.PlatformID_platformIDKick, nil
	default:
		err = fmt.Errorf("unknown platform ID %q", s)
	}
	return
}
