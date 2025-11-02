
all: go/streamcontrol_grpc

go/streamcontrol_grpc:
	protoc --go_out=. --go-grpc_out=. chatwebhook.proto
	rm -rf go
	mkdir -p go
	mv github.com/xaionaro-go/chatwebhook/pkg/grpc/protobuf/go/chatwebhook_grpc go/chatwebhook_grpc
	rm -rf github.com
