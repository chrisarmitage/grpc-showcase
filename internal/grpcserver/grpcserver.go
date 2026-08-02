package grpcserver

import (
	"context"
	"log"

	pb "github.com/chrisarmitage/grpc-showcase/proto"
	"google.golang.org/protobuf/types/known/emptypb"
)

type GrpcServer struct {
	pb.UnimplementedUniServiceServer
}

func (s *GrpcServer) Status(ctx context.Context, req *emptypb.Empty) (*pb.StatusResponse, error) {
	log.Printf("Received request for status request")
	return &pb.StatusResponse{Msg: "healthy"}, nil
}
