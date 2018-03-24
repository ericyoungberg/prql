package cmd

import (
  "os"
  "fmt"
  "time"
  "strings"
  "strconv"

  "github.com/spf13/cobra"
  "github.com/prql/prql/util"
  log "github.com/sirupsen/logrus"
  "github.com/olekukonko/tablewriter"
)


type Params struct{
  quiet  bool
  living bool

  username string
  host     string
  database string
  origins  string
}

const (
  tokenFile string = "/var/lib/prql/tokens"
)

var (
  params Params
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
    if params.username == "" {
      log.Fatal("Missing username [-u]")
    } else if params.host == "" {
      log.Fatal("Missing host [-H]")
    } else if params.database == "" {
      log.Fatal("Missing database [-d]")
    }

    timeSeed := strconv.Itoa(int(time.Now().Unix()))
    token    := util.CreateHash(strings.Join([]string{params.username, params.host, params.database, timeSeed}, ""))  
    password := util.GetPassword()

    entry := []string{token, params.username, password, params.host, params.database, params.origins, strconv.FormatBool(params.living)}

    err := util.AppendEntry(tokenFile, entry)
    if err != nil {
      log.Fatal("Could not generate new token.") 
      log.Fatal(err) 
    }
  },
}


var listTokensCmd = &cobra.Command{
  Use: "list",
  Short: "List all available tokens",
  Run: func(cmd *cobra.Command, args []string) {
    entries := util.ParseEntryFile(tokenFile)

    if params.quiet {
      tokens := make([]string, len(entries))

      for i, entry := range entries {
        tokens[i] = entry[0]
      }

      fmt.Println(strings.Join(tokens, " "))
    } else {
      table := tablewriter.NewWriter(os.Stdout)
      table.SetHeader([]string{"Token", "Username", "Host Name", "Database", "Origins", "Living"})

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
    entries := util.ParseEntryFile(tokenFile)

    for _, token := range args {
      for i, entry := range entries {
        if entry[0] == token {
          entries = append(entries[:i], entries[i+1:]...)

          fmt.Printf("Deleting %s\n", token) 
          break
        }
      }
    }

    err := util.WriteEntryFile(tokenFile, entries)
    if err != nil {
      log.Error("Could not write changes to tokens file")
      log.Error(err)
    }
  },
}


func init() {
  newTokenCmd.Flags().StringVarP(&params.username, "user", "u", "", "Database user associated with the new token")
  newTokenCmd.Flags().StringVarP(&params.host, "host", "H", "", "Database host associated with the new token. Must be a valid database host name defined using the databases command.")
  newTokenCmd.Flags().StringVarP(&params.database, "database", "d", "", "Database associated with the new token.")
  newTokenCmd.Flags().StringVarP(&params.origins, "origins", "", "", "Comma-delimited list of origins that are allowed to use the token.")
  newTokenCmd.Flags().BoolVarP(&params.living, "living", "l", false, "Keep connection alive, regardless of token usage frequency.")

  listTokensCmd.Flags().BoolVarP(&params.quiet, "quiet", "q", false, "Only display tokens")

  tokensCmd.AddCommand(newTokenCmd)
  tokensCmd.AddCommand(listTokensCmd)
  tokensCmd.AddCommand(removeTokenCmd)

  rootCmd.AddCommand(tokensCmd)
}
