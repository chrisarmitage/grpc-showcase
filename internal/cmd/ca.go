package cmd

import (
	"github.com/chrisarmitage/grpc-showcase/internal/cmd/ca"
	"github.com/spf13/cobra"
)

const (
	caFuncName = "ca"
	caCmdDesc  = "CA commands"
)

var caCmd = &cobra.Command{
	Use:   caFuncName,
	Short: caCmdDesc,
}

func CaCmd() *cobra.Command {
	caCmd.AddCommand(ca.GenerateCmd())

	return caCmd
}
