package main

import (
  "github.com/prql/prql/lib"
  log "github.com/sirupsen/logrus"
)


func main() {
  config, err := lib.GetConfig()
  if err != nil {
    log.Fatal("could not open prql.toml")     
  }

  populateDatabasePool()
  populateTokenPool()

  defer closeDatabaseConnections()

  startServer(&config)
}
