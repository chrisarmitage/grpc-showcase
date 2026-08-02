package grpcserver

import (
	"context"
	"fmt"
	"log"
	"time"

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

func (s *GrpcServer) StatusStream(req *emptypb.Empty, stream pb.UniService_StatusStreamServer) error {
	log.Printf("Received request for status stream")
	for i := range 5 {
		if err := stream.Send(&pb.StatusResponse{Msg: fmt.Sprintf("healthy %d", i)}); err != nil {
			return err
		}
		time.Sleep(1 * time.Second)
	}
	return nil
}
