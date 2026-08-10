package main

import (
	"os"

	"github.com/chrisarmitage/grpc-showcase/internal/cmd"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "grpc-showcase",
	Short: "A gRPC showcase application",
}

func main() {

	rootCmd.AddCommand(cmd.ServerCmd())
	rootCmd.AddCommand(cmd.ClientCmd())
	rootCmd.AddCommand(cmd.CaCmd())

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
