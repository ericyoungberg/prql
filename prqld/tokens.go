package main

import (
  "github.com/prql/prql/lib/pools"
)


var tokenPool map[string]pools.TokenEntry

func populateTokenPool() {
  tokenPool = make(map[string]pools.TokenEntry)
  tokenPool = pools.GetTokenEntries()
}
