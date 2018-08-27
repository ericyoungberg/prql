package cmd

import (
  "github.com/spf13/cobra"
  "github.com/prql/prql/lib"
  log "github.com/sirupsen/logrus"
)

var statusCmd = &cobra.Command{
  Use: "status",
  Short: "Check the status of prqld",
  Run: func(cmd *cobra.Command, args []string) {
    lib.CheckPrQLd()
    log.Info("status") 
  },
}

func init() {
  rootCmd.AddCommand(statusCmd)
}
