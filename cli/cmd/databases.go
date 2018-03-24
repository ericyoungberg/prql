package cmd

import (
  "fmt"
  "os"

  "github.com/spf13/cobra"
  "github.com/prql/prql/util"
  "github.com/olekukonko/tablewriter"
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
    fmt.Println(cmd.Short)
  },
}


var listDatabasesCmd = &cobra.Command{
  Use: "list",
  Short: "List all available databases",
  Run: func(cmd *cobra.Command, args []string) {
    table := tablewriter.NewWriter(os.Stdout)
    table.SetHeader([]string{"Host Name", "Driver", "Host", "Port", "SSL"})

    entries := util.ParseEntryFile("/var/lib/prql/databases")
    table.AppendBulk(entries)
    table.Render()
  },
}


var removeDatabaseCmd = &cobra.Command{
  Use: "remove",
  Short: "Remove database location from system. This action is permanent.",
  Run: func(cmd *cobra.Command, args []string) {
    fmt.Println(cmd.Short) 
  },
}


func init() {
  databasesCmd.AddCommand(newDatabaseCmd)
  databasesCmd.AddCommand(listDatabasesCmd)
  databasesCmd.AddCommand(removeDatabaseCmd)

  rootCmd.AddCommand(databasesCmd)
}
