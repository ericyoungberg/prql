package cmd

import (
  "os"
  "fmt"
  "strings"
  "strconv"

  "github.com/spf13/cobra"
  "github.com/prql/prql/lib"
  log "github.com/sirupsen/logrus"
  "github.com/olekukonko/tablewriter"
)

type DatabaseParams struct{
  quiet bool
  ssl   bool

  hostName string
  driver   string
  host     string

  port int
}

const (
  databaseFile string = "/var/lib/prql/databases"
)

var (
  databaseParams DatabaseParams
)


var databasesCmd = &cobra.Command{
  Use: "databases",
  Short: "Add, delete, or view all databases added to the system",
  Run: func(cmd *cobra.Command, args []string) {
    fmt.Println(cmd.Short) 
  },
}


var newDatabaseCmd = &cobra.Command{
  Use: "new",
  Short: "Add a database to be referenced by the system",
  Run: func(cmd *cobra.Command, args []string) {
    if databaseParams.hostName == "" {
      log.Fatal("Missing host name [-n]")
    } else if databaseParams.driver == "" {
      log.Fatal("Missing driver [-d]")
    } else if databaseParams.port == 0 {
      log.Fatal("Missing port [-p]")
    }

    entry := []string{databaseParams.hostName, databaseParams.driver, databaseParams.host, strconv.Itoa(databaseParams.port), strconv.FormatBool(databaseParams.ssl)}

    err := lib.AppendEntry(databaseFile, entry)
    if err != nil {
      log.Error("Could not add database") 
      log.Fatal(err) 
    }

    fmt.Printf("Added database %s\n", databaseParams.hostName)

    refreshServerPool("databases")
  },
}


var listDatabasesCmd = &cobra.Command{
  Use: "list",
  Short: "List all available databases",
  Run: func(cmd *cobra.Command, args []string) {
    entries := lib.ParseEntryFile(databaseFile)

    if databaseParams.quiet {
      names := make([]string, len(entries)) 

      for i, entry := range entries {
        names[i] = entry[0]
      }

      fmt.Println(strings.Join(names, " "))
    } else {
      table := tablewriter.NewWriter(os.Stdout)
      table.SetHeader([]string{"Host Name", "Driver", "Host", "Port", "SSL"})

      table.AppendBulk(entries)
      table.Render()
    }
  },
}


var removeDatabaseCmd = &cobra.Command{
  Use: "remove [names]",
  Short: "Remove database location from system. This action is permanent.",
  Run: func(cmd *cobra.Command, args []string) {
    entries := lib.ParseEntryFile(databaseFile)
    entries = lib.RemoveByColumn(args, entries, 0)

    err := lib.WriteEntryFile(databaseFile, entries)
    if err != nil {
      log.Error("Could not write changes to tokens file")
      log.Error(err)
      return
    }

    refreshServerPool("databases")
  },
}


func init() {
  newDatabaseCmd.Flags().StringVarP(&databaseParams.hostName, "name", "n", "", "Host name used to reference this server from the tokens")
  newDatabaseCmd.Flags().StringVarP(&databaseParams.driver, "driver", "d", "", "Database type (postgresql, mysql)")
  newDatabaseCmd.Flags().StringVarP(&databaseParams.host, "host", "H", "0.0.0.0", "Location of the database server")
  newDatabaseCmd.Flags().IntVarP(&databaseParams.port, "port", "p", 0, "Port of the database server")

  listDatabasesCmd.Flags().BoolVarP(&databaseParams.quiet, "quiet", "q", false, "Only display host names")

  databasesCmd.AddCommand(newDatabaseCmd)
  databasesCmd.AddCommand(listDatabasesCmd)
  databasesCmd.AddCommand(removeDatabaseCmd)

  rootCmd.AddCommand(databasesCmd)
}
