package cmd

import (
  "os"
  "fmt"
  "strconv"
  "net/url"
  "net/http"

  "github.com/spf13/cobra"
  "github.com/prql/prql/lib"
  log "github.com/sirupsen/logrus"
)


var rootCmd = &cobra.Command{
  Use: "prql ",
  Short: "PrQL is a service for executing SQL queries over HTTP",
  Long: ``,
}


func Execute() {
  if err := rootCmd.Execute(); err != nil {
    fmt.Println(err)
    os.Exit(1)
  }
}


func refreshServerPool(poolName string) {
  config, err := lib.GetConfig()
  if err != nil {
    log.Fatal(err) 
  }

  updateErrMsg := "Could not refresh " + poolName + " pool in prqld"

  endpoint := url.URL{
    Scheme: "http", 
    Host: ("127.0.0.1:" + strconv.Itoa(config.Port)), 
    Path: ("refresh-" + poolName),
  }

  req, err := http.NewRequest("GET", endpoint.String(), nil)
  if err != nil {
    log.Fatal(err)
  }

  req.Header.Set(config.Headers.Secret, config.Secret)

  client := &http.Client{}
  res, err := client.Do(req)
  if err != nil {
    log.Fatal(updateErrMsg) 
  }
  res.Body.Close();

  if res.StatusCode != http.StatusOK {
    log.Error(updateErrMsg)
  }
}
