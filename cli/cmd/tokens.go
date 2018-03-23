package cmd

import (
  "fmt"
  "os"

  "github.com/prql/prql/util"
  "github.com/spf13/cobra"
  "github.com/olekukonko/tablewriter"
)


var tokensCmd = &cobra.Command{
  Use: "tokens",
  Short: "Generate, delete, or view all PrQL tokens",
  Run: func(cmd *cobra.Command, args []string) {
    fmt.Println(cmd.Short) 
  },
}


var newTokenCmd = &cobra.Command{
  Use: "new",
  Short: "Generate a new PrQL token for the given credentials",
  Run: func(cmd *cobra.Command, args []string) {
    fmt.Println(cmd.Short)
  },
}


var listTokensCmd = &cobra.Command{
  Use: "list",
  Short: "List all available tokens",
  Run: func(cmd *cobra.Command, args []string) {
    table := tablewriter.NewWriter(os.Stdout)
    table.SetHeader([]string{"Token", "Username", "Server", "Database", "Domains", "Living"})

    entries := util.ParseEntryFile("/var/lib/prql/tokens")
    for _, entry := range entries {
      table.Append(append(entry[:2], entry[3:]...))
    }

    table.Render()
  },
}


var removeTokenCmd = &cobra.Command{
  Use: "remove",
  Short: "Remove token. This action is permanent.",
  Run: func(cmd *cobra.Command, args []string) {
    fmt.Println(cmd.Short) 
  },
}


func init() {
  tokensCmd.AddCommand(newTokenCmd)
  tokensCmd.AddCommand(listTokensCmd)
  tokensCmd.AddCommand(removeTokenCmd)

  rootCmd.AddCommand(tokensCmd)
}
