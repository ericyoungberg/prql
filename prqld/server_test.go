package main

import (
  "http"
  "testing"

  "github.com/prql/prql/lib/pools"
)


type AuthTest struct {
  Entry pools.TokenEntry
  Origin string
  Outcome bool
}

type RequestTest struct {
  Request http.Request
  Token string
  Query string 
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

  requestTests = []RequestTest{
    RequestTest{
      
    },

    RequestTest{
    
    },

    RequestTest{
    
    },

    RequestTest{
    
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

func TestReadToken(t *testing.T) {
  for _, test := range requestTests {
    if token := readToken(&test.Request); token != test.Token {
      t.Errorf("expected %s, got %t", test.Token, token) 
    }
  }
}

func TestReadQuery(t *testing.T) {
  for _, test := range requestTests {
    if query := readQuery(&test.Query); query != test.Query {
      t.Errorf("expected %s, got %t", test.Query, query) 
    }
  }
}
