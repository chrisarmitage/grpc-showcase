package server

import (
	"log"
	"net"

	pb "github.com/chrisarmitage/grpc-showcase/proto"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"github.com/chrisarmitage/grpc-showcase/internal/grpcserver"
)

const (
	startFuncName = "start"
	startCmdDesc  = "Start the gRPC server"
	modeInsecure  = "insecure"
	modeTls       = "tls"
)

var tlsMode string

func StartCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   startFuncName,
		Short: startCmdDesc,
		Run:   runServerStart,
	}

	cmd.Flags().StringVar(&tlsMode, "tls", modeInsecure, "TLS mode (insecure|tls)")

	return cmd
}

func runServerStart(cmd *cobra.Command, args []string) {
	lis, err := net.Listen("tcp", ":5050")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	var serverOptions []grpc.ServerOption

	switch tlsMode {
	case modeInsecure:
		log.Printf("Starting gRPC server in insecure mode")
	case modeTls:
		creds, err := credentials.NewServerTLSFromFile("pki/server.crt", "pki/server.key")
		if err != nil {
			log.Fatalf("Failed to load TLS credentials from pki/server.crt/pki/server.key: %v", err)
		}
		serverOptions = append(serverOptions, grpc.Creds(creds))
		log.Printf("Starting gRPC server in TLS mode")
	default:
		log.Fatalf("Invalid TLS mode %q. Accepted values are: %s|%s", tlsMode, modeInsecure, modeTls)
	}

	s := grpc.NewServer(serverOptions...)
	pb.RegisterUniServiceServer(s, &grpcserver.GrpcServer{})

	log.Printf("gRPC server listening on port 5050")
	log.Printf("Use Ctrl+C to stop the server")

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
