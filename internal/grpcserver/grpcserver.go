package grpcserver

import (
	"context"
	"fmt"
	"io"
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

func (s *GrpcServer) Bidir(stream pb.UniService_BidirServer) error {
	log.Printf("Received request for bidir stream")

	clientDone := make(chan error, 1)
	go func() {
		// listen loop
		for {

			msg, err := stream.Recv()
			if err == io.EOF {
				log.Printf("Bidir stream completed, closed by client")
				clientDone <- nil
				return
			}
			if err != nil {
				clientDone <- err
				return
			}

			log.Printf("Bidir <-  : %s", msg.Msg)

		}

	}()

	// send a message every 500 milliseconds until the client closes the stream
	counter := 0
	for {
		select {
		case err := <-clientDone:
			return err
		default:
			if err := stream.Send(&pb.BidirMessage{Msg: fmt.Sprintf("server message %d", counter)}); err != nil {
				return err
			}
			log.Printf("Bidir  -> : server message %d", counter)
			counter++
			time.Sleep(500 * time.Millisecond)
		}
	}
}
