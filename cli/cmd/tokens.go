package cmd

import (
  "os"
  "fmt"
  "time"
  "strings"
  "strconv"

  "github.com/spf13/cobra"
  "github.com/prql/prql/lib"
  log "github.com/sirupsen/logrus"
  "github.com/olekukonko/tablewriter"
)


type TokenParams struct{
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
  tokenParams TokenParams
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
    _, err := lib.GetConfig()
    if err != nil {
      log.Fatal("Could not load configuration", err) 
    }

    if tokenParams.username == "" {
      log.Fatal("Missing username [-u]")
    } else if tokenParams.host == "" {
      log.Fatal("Missing host [-H]")
    } else if tokenParams.database == "" {
      log.Fatal("Missing database [-d]")
    }

    timeSeed := strconv.Itoa(int(time.Now().Unix()))
    token    := lib.CreateHash(strings.Join([]string{tokenParams.username, tokenParams.host, tokenParams.database, timeSeed}, ""))  
    password, err := lib.GetPassword(tokenParams.username)
    if err != nil {
      log.Fatal(err) 
      return
    }

    entry := []string{token, tokenParams.username, password, tokenParams.host, tokenParams.database, tokenParams.origins, strconv.FormatBool(tokenParams.living)}

    err = lib.AppendEntry(tokenFile, entry)
    if err != nil {
      log.Fatal("Could not generate new token.") 
      log.Fatal(err) 
      return
    }

    fmt.Printf("Generated Token %s\n", token)

    refreshServerPool("tokens")
  },
}


var listTokensCmd = &cobra.Command{
  Use: "list",
  Short: "List all available tokens",
  Run: func(cmd *cobra.Command, args []string) {
    entries := lib.ParseEntryFile(tokenFile)

    config, err := lib.GetConfig()
    if err != nil {
      log.Fatal("Could not load configuration", err) 
    }

    fmt.Printf("%+v\n", config)


    if tokenParams.quiet {
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
    entries := lib.ParseEntryFile(tokenFile)
    entries = lib.RemoveByColumn(args, entries, 0)

    err := lib.WriteEntryFile(tokenFile, entries)
    if err != nil {
      log.Error("Could not write changes to tokens file")
      log.Error(err)
    }

    refreshServerPool("tokens")
  },
}


func init() {
  newTokenCmd.Flags().StringVarP(&tokenParams.username, "user", "u", "", "Database user associated with the new token")
  newTokenCmd.Flags().StringVarP(&tokenParams.host, "host", "H", "", "Database host associated with the new token. Must be a valid database host name defined using the databases command.")
  newTokenCmd.Flags().StringVarP(&tokenParams.database, "database", "d", "", "Database associated with the new token.")
  newTokenCmd.Flags().StringVarP(&tokenParams.origins, "origins", "o", "", "Comma-delimited list of origins that are allowed to use the token.")
  newTokenCmd.Flags().BoolVarP(&tokenParams.living, "living", "l", false, "Keep connection alive, regardless of token usage frequency.")

  listTokensCmd.Flags().BoolVarP(&tokenParams.quiet, "quiet", "q", false, "Only display tokens")

  tokensCmd.AddCommand(newTokenCmd)
  tokensCmd.AddCommand(listTokensCmd)
  tokensCmd.AddCommand(removeTokenCmd)

  rootCmd.AddCommand(tokensCmd)
}
