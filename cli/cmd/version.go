package cmd

import (
  "fmt"

  "github.com/prql/prql/cli/version"
  "github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
  Use: "version, v",
  Aliases: []string{"v"},
  Short: "Check the version of prql and prqld",
  Run: func(_ *cobra.Command, _ []string) {
    fmt.Printf(`PrQL Versions:

  prql  : %s
  prlqd : %s

`, version.VERSION, "v0.1.1")
  },
}

func init() {
  rootCmd.AddCommand(versionCmd)
}
