package cmd

import (
	"github.com/chrisarmitage/grpc-showcase/internal/cmd/server"
	"github.com/spf13/cobra"
)

const (
	serverFuncName = "server"
	serverCmdDesc  = "Server commands"
)

var serverCmd = &cobra.Command{
	Use:   serverFuncName,
	Short: serverCmdDesc,
}

func ServerCmd() *cobra.Command {
	serverCmd.AddCommand(server.StartCmd())

	return serverCmd
}
