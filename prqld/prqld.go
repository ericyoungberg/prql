package main

import (
  "github.com/prql/prql/lib"
)


func main() {
  PopulateDatabasePool(false)
  PopulateTokenPool(false)

  defer CloseDatabaseConnections()

  StartServer(&lib.Config{Port: 1999})
}
