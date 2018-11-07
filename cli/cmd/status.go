package cmd

import (
  "github.com/spf13/cobra"
  "github.com/prql/prql/lib"
)

var statusCmd = &cobra.Command{
  Use: "status",
  Short: "Check the status of prqld",
  Run: func(cmd *cobra.Command, args []string) {
    lib.CheckServerStatus()
  },
}

func init() {
  rootCmd.AddCommand(statusCmd)
}
