package main

import (
  "testing"
  "github.com/prql/prql/lib/pools"
)


type AuthTest struct {
  Entry pools.TokenEntry
  Origin string
  Outcome bool
}


var (
  tokenEntry = pools.TokenEntry{
    Origins: []string{
      "place.com",
      "otherplace.com",
      "127.0.0.1",
    },
  }

  emptyTokenEntry = pools.TokenEntry{}

  authTests = []AuthTest{
    AuthTest{
      Entry: tokenEntry,
      Origin: "otherplace.com",
      Outcome: true,
    },

    AuthTest{
      Entry: tokenEntry,
      Origin: "someplaceelse.com",
      Outcome: false,
    },

    AuthTest{
      Entry: tokenEntry,
      Origin: "",
      Outcome: false,
    },

    AuthTest{
      Entry: emptyTokenEntry,
      Origin: "",
      Outcome: true,
    },

    AuthTest{
      Entry: emptyTokenEntry,
      Origin: "place.com",
      Outcome: true,
    },
  }
)


func TestAuthorizedOrigin(t *testing.T) {
  for _, test := range authTests {
    if result := authorizedOrigin(test.Origin, test.Entry); result != test.Outcome {
      t.Errorf("expected %t, got %t [value=%s] [origins=%v]", test.Outcome, result, test.Origin, test.Entry.Origins) 
    }
  }
}
