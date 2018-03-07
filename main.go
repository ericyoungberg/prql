package main

func main() {
  PopulateDatabasePool()
  PopulateTokenPool()
  StartServer(&Config{Port: 1999})
}
