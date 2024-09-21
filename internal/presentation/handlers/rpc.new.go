package handlers

import "github.com/sheginabo/go-quick-grpc/internal/pb"

type GrpcApi struct {
	pb.UnimplementedGoQuickGRPCServer
}

func NewGrpcApi() *GrpcApi {
	return &GrpcApi{}
}
