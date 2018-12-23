package cmd

import (
  "fmt"
  "net/http"

  "github.com/spf13/cobra"
  "github.com/prql/prql/lib"
  log "github.com/sirupsen/logrus"
  "github.com/prql/prql/cli/version"
)

var versionCmd = &cobra.Command{
  Use: "version",
  Aliases: []string{"v"},
  Short: "Check the version of prql and prqld",
  Run: func(_ *cobra.Command, _ []string) {
    config, err := lib.GetConfig() 
    if err != nil {
      log.Fatal(err) 
    }

    var dVersion string
    dVersionEndpoint := fmt.Sprintf("127.0.0.1:%d/version", config.Port)
    res, err := http.Get(dVersionEndpoint) 
    if err != nil {
      dVersion = "Unavailable" 
    } else {
      dVersion = ""
    }

    fmt.Println(res, dVersion)

    fmt.Printf(`PrQL Versions:

  prql  : %s
  prlqd : %s

`, version.VERSION, dVersion)
  },
}

func init() {
  rootCmd.AddCommand(versionCmd)
}
