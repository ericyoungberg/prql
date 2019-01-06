package main

import (
  "github.com/prql/prql/lib"
  "github.com/sirupsen/logrus"
  "github.com/prql/prql/prqld/databases"
)

func main() {
  config, err := lib.GetConfig()
  if err != nil {
    logrus.Fatal("could not open prql.toml")     
  }

  setupLogger(&config)

  databases.populateTokenPool()
  databases.populateDatabasePool()

  defer databases.closeDatabaseConnections()

  server := &Server{}
  server.StartFromConfig(&config)
}
