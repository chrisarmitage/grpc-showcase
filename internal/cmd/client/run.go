package client

import (
	"context"
	"log"

	"github.com/chrisarmitage/grpc-showcase/proto"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"
)

const (
	runFuncName = "run"
	runCmdDesc  = "Run the gRPC client"
)

func RunCmd() *cobra.Command {
	return &cobra.Command{
		Use:   runFuncName,
		Short: runCmdDesc,
		Run:   runClient,
	}
}

func runClient(cmd *cobra.Command, args []string) {
	conn, err := grpc.NewClient("localhost:5050", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	client := proto.NewStatusServiceClient(conn)
	resp, err := client.Status(context.Background(), &emptypb.Empty{})
	if err != nil {
		log.Fatalf("Error calling Status: %v", err)
	}

	log.Printf("Status response: %s", resp.Msg)
}
