package cmd

import (
	"github.com/chrisarmitage/grpc-showcase/internal/cmd/client"
	"github.com/spf13/cobra"
)

const (
	clientFuncName = "client"
	clientCmdDesc  = "Client commands"
)

var clientCmd = &cobra.Command{
	Use:   clientFuncName,
	Short: clientCmdDesc,
}

func ClientCmd() *cobra.Command {
	clientCmd.AddCommand(client.RunCmd())

	return clientCmd
}
