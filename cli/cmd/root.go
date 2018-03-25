package cmd

import (
  "os"
  "fmt"

  "github.com/spf13/cobra"
)


var params interface{}
  

var rootCmd = &cobra.Command{
  Use: "prql ",
  Short: "PrQL is a service for executing SQL queries over HTTP",
  Long: ``,
  Run: func(cmd *cobra.Command, args []string) {
    fmt.Println("From within the CLI") 
  },
}


func Execute() {
  if err := rootCmd.Execute(); err != nil {
    fmt.Println(err)
    os.Exit(1)
  }
}
