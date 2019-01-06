package tokens

import (
  "github.com/prql/prql/lib"
)


var TokenPool map[string]lib.TokenEntry

func populateTokenPool() {
  tokenPool = make(map[string]lib.TokenEntry)
  tokenPool = lib.GetTokenEntries()
}
