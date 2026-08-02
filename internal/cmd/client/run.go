package client

import (
	"context"
	"io"
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
	modeSingle  = "single"
	modeStream  = "stream"
)

var runMode string

func RunCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   runFuncName,
		Short: runCmdDesc,
		Run:   runClient,
	}

	cmd.Flags().StringVar(&runMode, "mode", modeSingle, "Client run mode (single|stream)")

	return cmd
}

func runClient(cmd *cobra.Command, args []string) {
	conn, err := grpc.NewClient("localhost:5050", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	client := proto.NewUniServiceClient(conn)

	switch runMode {
	case modeSingle:
		resp, err := client.Status(context.Background(), &emptypb.Empty{})
		if err != nil {
			log.Fatalf("Error calling Status: %v", err)
		}

		log.Printf("Status response: %s", resp.Msg)
	case modeStream:
		stream, err := client.StatusStream(context.Background(), &emptypb.Empty{})
		if err != nil {
			log.Fatalf("Error calling StatusStream: %v", err)
		}

		for {
			resp, err := stream.Recv()
			if err == io.EOF {
				log.Printf("Status stream completed, closed by server")
				break
			}
			if err != nil {
				log.Fatalf("Error receiving StatusStream response: %v", err)
			}

			log.Printf("Status stream response: %s", resp.Msg)
		}
	default:
		log.Fatalf("Invalid mode %q. Accepted values are: %s|%s", runMode, modeSingle, modeStream)
	}
}
