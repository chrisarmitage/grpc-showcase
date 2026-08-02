package server

import (
	"log"
	"net"

	pb "github.com/chrisarmitage/grpc-showcase/proto"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"

	"github.com/chrisarmitage/grpc-showcase/internal/grpcserver"
)

const (
	startFuncName = "start"
	startCmdDesc  = "Start the gRPC server"
)

func StartCmd() *cobra.Command {
	return &cobra.Command{
		Use:   startFuncName,
		Short: startCmdDesc,
		Run:   runServerStart,
	}
}

func runServerStart(cmd *cobra.Command, args []string) {
	lis, err := net.Listen("tcp", ":5050")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterUniServiceServer(s, &grpcserver.GrpcServer{})

	log.Printf("gRPC server listening on port 5050")
	log.Printf("Use Ctrl+C to stop the server")

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
