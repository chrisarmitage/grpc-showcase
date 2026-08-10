package client

import (
	"context"
	"io"
	"log"
	"time"

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
	modeBidir   = "bidir"
)

var runMode string

func RunCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   runFuncName,
		Short: runCmdDesc,
		Run:   runClient,
	}

	cmd.Flags().StringVar(&runMode, "mode", modeSingle, "Client run mode (single|stream|bidir)")

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
	case modeBidir:
		stream, err := client.Bidir(context.Background())
		if err != nil {
			log.Fatalf("Error calling Bidir: %v", err)
		}

		recvDone := make(chan error, 1)
		go func() {
			// listen loop
			for {
				resp, err := stream.Recv()
				if err == io.EOF {
					recvDone <- nil
					return
				}
				if err != nil {
					recvDone <- err
					return
				}

				log.Printf("Bidir <-  : %s", resp.Msg)
			}
		}()

		messages := []string{"client-1", "client-2", "client-3"}
		for _, message := range messages {
			log.Printf("Bidir  -> : %s", message)
			if err := stream.Send(&proto.BidirMessage{Msg: message}); err != nil {
				log.Fatalf("Error sending Bidir message: %v", err)
			}

			time.Sleep(2 * time.Second)
		}

		if err := stream.CloseSend(); err != nil {
			log.Fatalf("Error closing Bidir stream: %v", err)
		}

		if err := <-recvDone; err != nil {
			log.Fatalf("Error receiving Bidir response: %v", err)
		}
	default:
		log.Fatalf("Invalid mode %q. Accepted values are: %s|%s|%s", runMode, modeSingle, modeStream, modeBidir)
	}
}
