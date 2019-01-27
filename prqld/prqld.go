package main

import (
  "github.com/prql/prql/lib"
  "github.com/sirupsen/logrus"
)

func main() {
  config, err := lib.GetConfig()
  if err != nil {
    logrus.Fatal("could not open prql.toml")     
  }

  setupLogger(&config)

  defer closeDatabaseConnections()

  server := &Server{}
  server.StartFromConfig(&config)
}
