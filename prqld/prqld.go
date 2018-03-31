package main

func main() {
  PopulateDatabasePool(false)
  PopulateTokenPool(false)

  defer CloseDatabaseConnections()

  StartServer(&Config{Port: 1999})
}
