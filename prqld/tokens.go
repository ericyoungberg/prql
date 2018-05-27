package main

import (
  "github.com/prql/prql/lib"
)


var tokenPool map[string]lib.TokenEntry

func populateTokenPool() {
  tokenPool = make(map[string]lib.TokenEntry)
  tokenPool = lib.GetTokenEntries()
}
