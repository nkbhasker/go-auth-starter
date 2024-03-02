package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use: "authenticator",
	RunE: func(cmd *cobra.Command, _args []string) error {
		return SrvStart(cmd.Version)
	},
}

func Version(version string) {
	rootCmd.Version = version
}

func Execute() error {
	rootCmd.AddCommand(schemaCmd, srvStartCmd)

	return rootCmd.Execute()
}
