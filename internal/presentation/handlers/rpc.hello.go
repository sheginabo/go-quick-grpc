package handlers

import (
	"context"
	"github.com/sheginabo/go-quick-grpc/internal/pb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

func (c *GrpcApi) SendHello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloResponse, error) {
	rsp := &pb.HelloResponse{
		Message:        "Hello " + req.Message,
		Timestamp:      timestamppb.New(time.Now()),
		TimestampMilli: time.Now().UnixMilli(),
	}
	return rsp, nil
}
