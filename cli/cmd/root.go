package cmd

import (
  "os"
  "fmt"
  "net/url"
  "net/http"

  "github.com/spf13/cobra"
  log "github.com/sirupsen/logrus"
)


var params interface{}
  

var rootCmd = &cobra.Command{
  Use: "prql ",
  Short: "PrQL is a service for executing SQL queries over HTTP",
  Long: ``,
  Run: func(cmd *cobra.Command, args []string) {
    fmt.Println("From within the CLI") 
  },
}


func Execute() {
  if err := rootCmd.Execute(); err != nil {
    fmt.Println(err)
    os.Exit(1)
  }
}

func refreshServerPool(poolName string) {
  endpoint := url.URL{Scheme: "http", Host: "127.0.0.1:1999", Path: "refresh-" + poolName}
  res, err := http.Get(endpoint.String())
  if err != nil {
    log.Fatal(err) 
  }

  res.Body.Close();

  if res.StatusCode != http.StatusOK {
    log.Error("Could not refresh " + poolName + " pool in prqld")
  }
}
