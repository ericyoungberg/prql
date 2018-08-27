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

var (
  tokenParams TokenParams
)


var tokensCmd = &cobra.Command{
  Use: "tokens",
  Short: "Generate, delete, or view all PrQL tokens",
}


var newTokenCmd = &cobra.Command{
  Use: "new",
  Short: "Generate a new PrQL token for the given credentials",
  Run: func(cmd *cobra.Command, args []string) {
    if tokenParams.username == "" {
      log.Fatal("Missing username [-u]")
    } else if tokenParams.host == "" {
      log.Fatal("Missing host [-H]")
    } else if tokenParams.database == "" {
      log.Fatal("Missing database [-d]")
    }

    origins := tokenParams.origins
    if origins != "" {
      originEntries := strings.Split(origins, ",")
      validOrigins := make([]string, 0, len(originEntries))
      for _, origin := range originEntries {
        println(fmt.Sprintf("origin: |%s| |%s|", origin, strings.Replace(origin, " ", "", -1)))
        if strings.Replace(origin, " ", "", -1) != "" {
          validOrigins = append(validOrigins, origin)
        }
      }
    }

    println(fmt.Sprintf("Origins: |%s|", origins))

    timeSeed := strconv.Itoa(int(time.Now().Unix()))
    token    := lib.CreateHash(strings.Join([]string{tokenParams.username, tokenParams.host, tokenParams.database, timeSeed}, ""))  
    password, err := lib.GetPassword(tokenParams.username)
    if err != nil {
      log.Fatal(err) 
      return
    }

    entry := []string{token, tokenParams.username, password, tokenParams.host, tokenParams.database, origins, strconv.FormatBool(tokenParams.living)}

    err = lib.AppendEntry(lib.Sys.TokenFile, entry)
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
    entries := lib.ParseEntryFile(lib.Sys.TokenFile)

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
    entries := lib.ParseEntryFile(lib.Sys.TokenFile)
    entries = lib.RemoveByColumn(args, entries, 0)

    err := lib.WriteEntryFile(lib.Sys.TokenFile, entries)
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
