package main

import (
  "github.com/prql/prql/lib"
  "github.com/sirupsen/logrus"
)

// Populated by ldflags
var (
  VERSION string
)

func main() {
  config, err := lib.GetConfig()
  if err != nil {
    logrus.Fatal("could not open prql.toml")     
  }

  setupLogger(&config)

  populateTokenPool()
  populateDatabasePool()

  defer closeDatabaseConnections()

  startServer(&config)
}
