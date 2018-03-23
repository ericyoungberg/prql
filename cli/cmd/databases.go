package cmd

import (
  "fmt"

  "github.com/spf13/cobra"
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
    fmt.Println(cmd.Short)
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
