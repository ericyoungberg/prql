package cmd

import (
  "os"
  "fmt"
  "strings"

  "github.com/spf13/cobra"
  "github.com/prql/prql/util"
  "github.com/olekukonko/tablewriter"
)

var (
  quietMode bool
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
    entries := util.ParseEntryFile("/var/lib/prql/tokens")

    if quietMode {
      tokens := make([]string, len(entries))

      for i, entry := range entries {
        tokens[i] = entry[0]
      }

      fmt.Println(strings.Join(tokens, " "))
    } else {
      table := tablewriter.NewWriter(os.Stdout)
      table.SetHeader([]string{"Token", "Username", "Server Name", "Database", "Domains", "Living"})

      for _, entry := range entries {
        table.Append(append(entry[:2], entry[3:]...))
      }

      table.Render()
    }
  },
}


var removeTokenCmd = &cobra.Command{
  Use: "remove [tokens]",
  Short: "Remove token. This action is permanent.",
  Run: func(cmd *cobra.Command, args []string) {
    entries := util.ParseEntryFile("/var/lib/prql/tokens")

    for _, token := range args {
      for i, entry := range entries {
        if entry[0] == token {
          entries = append(entries[:i], entries[i+1:]...)

          fmt.Printf("Deleting %s\n", token) 
          break
        }
      }
    }

    /*
    err := util.WriteEntryFile(entries)
    if err != nil {
    
    }
    */
  },
}


func init() {
  listTokensCmd.Flags().BoolVarP(&quietMode, "quiet", "q", false, "Only display tokens")

  tokensCmd.AddCommand(newTokenCmd)
  tokensCmd.AddCommand(listTokensCmd)
  tokensCmd.AddCommand(removeTokenCmd)

  rootCmd.AddCommand(tokensCmd)
}
