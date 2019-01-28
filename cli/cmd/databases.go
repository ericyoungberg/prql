package cmd

import (
  "os"
  "fmt"
  "strings"
  "strconv"

  "github.com/spf13/cobra"
  "github.com/prql/prql/lib"
  "github.com/prql/prql/lib/pools"
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

var (
  databaseParams DatabaseParams
  supportedDatabases = [...]string{ "postgres", "mysql" }
)

var databasesCmd = &cobra.Command{
  Use: "databases",
  Short: "Add, delete, or view all databases added to the system",
}

var listDatabasesCmd = &cobra.Command{
  Use: "list",
  Short: "List all available databases",
  Run: func(cmd *cobra.Command, args []string) {
    entries := pools.ParseEntryFile(lib.Sys.DatabaseFile)

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

    pool := pools.NewDatabasePool()

    dbSupported := false
    for _, supportedDatabase := range supportedDatabases {
      if databaseParams.driver == supportedDatabase {
        dbSupported = true 
        break
      }
    }
    if !dbSupported {
      log.Fatal(databaseParams.driver + " is not a supported driver. Supported drivers: " + strings.Join(supportedDatabases[:], ", ")) 
    }

    if _, nameUsed := pool.Entries[databaseParams.hostName]; nameUsed {
      log.Fatal("The host name " + databaseParams.hostName + " has already been used")
    }

    pool.AppendRecord([]string{
      databaseParams.hostName, 
      databaseParams.driver, 
      databaseParams.host, 
      strconv.Itoa(databaseParams.port), 
      strconv.FormatBool(databaseParams.ssl),
    })

    err := pool.Save()
    if err != nil {
      log.Error("Could not add database") 
      log.Fatal(err) 
    }
    fmt.Printf("Added database %s\n", databaseParams.hostName)

    refreshServerPool("databases")
  },
}

var removeDatabaseCmd = &cobra.Command{
  Use: "remove [names]",
  Short: "Remove database location from system. This action is permanent.",
  Run: func(cmd *cobra.Command, args []string) {
    pool := pools.NewDatabasePool()
    pool.Remove(args)

    err := pool.Save()
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
  newDatabaseCmd.Flags().StringVarP(&databaseParams.driver, "driver", "d", "", "Database type (" + strings.Join(supportedDatabases[:], ", ") + ")")
  newDatabaseCmd.Flags().StringVarP(&databaseParams.host, "host", "H", "0.0.0.0", "Location of the database server")
  newDatabaseCmd.Flags().IntVarP(&databaseParams.port, "port", "p", 0, "Port of the database server")

  listDatabasesCmd.Flags().BoolVarP(&databaseParams.quiet, "quiet", "q", false, "Only display host names")

  databasesCmd.AddCommand(newDatabaseCmd)
  databasesCmd.AddCommand(listDatabasesCmd)
  databasesCmd.AddCommand(removeDatabaseCmd)

  rootCmd.AddCommand(databasesCmd)
}
