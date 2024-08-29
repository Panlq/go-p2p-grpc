package cmd

import (
	"github.com/spf13/cobra"

	"github/panlq-github/go-p2p-grpc/internal/server"
)

var conf server.Config

func init() {
	conf = server.Config{}
}

func NewP2PCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "go-p2p",
		RunE: func(cmd *cobra.Command, args []string) error {
			return server.NewNode(conf).Start()
		},
	}

	cmd.PersistentFlags().StringVarP(&conf.NodeName, "node", "n", "node", "Node name")

	cmd.PersistentFlags().StringVarP(&conf.NodeAddr, "addr", "a", "127.0.0.1:30051", "Node listen address")

	cmd.PersistentFlags().StringVarP(&conf.ServiceDiscoveryAddress, "consul", "c", "127.0.0.1:8500", "Consul address")

	cmd.MarkFlagsRequiredTogether("node", "addr", "consul")

	return cmd
}
