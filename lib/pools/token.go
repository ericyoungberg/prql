package pools

import (
  "strconv"
  "strings"

  "github.com/prql/prql/lib"
  log "github.com/sirupsen/logrus"
)


type TokenEntry struct{
  Tag      string
  User     string
  Password string
  HostName string
  DBName   string

  Origins []string

  Living bool
}

type TokenPool struct {
  pool

  Entries map[string]TokenEntry
}

func (p *TokenPool) build() {
  tokens := make(map[string]TokenEntry) 

  for i, parts := range p.records {
    if len(parts) != 8 {
      log.Error("Invalid token entry at line " + strconv.Itoa(i + 1)) 
      continue
    }

    originEntries := strings.Split(parts[6], ",")
    origins := make([]string, len(originEntries))
    originIndex := 0
    for _, origin := range originEntries {
      if origin != "" {
        origins[originIndex] = origin
        originIndex += 1
      }
    }

    if len(origins) == 1 && origins[0] == "" {
      origins = nil 
    }

    living, err := strconv.ParseBool(parts[7])
    if err != nil {
      living = false 
    }

    password, err := lib.InsecureDecryptString(parts[3])
    if err != nil {
      log.Error("Couldn't decrypt password at line " + strconv.Itoa(i + 1)) 
      continue
    }

    tokens[parts[0]] = TokenEntry{
      Tag: parts[1]
      User: parts[2], 
      Password: password, 
      HostName: parts[4], 
      DBName: parts[5], 
      Origins: origins, 
      Living: living,
    }
  }

  p.Entries = tokens
}

func GetTokenPool() *TokenPool {
  if __TOKEN_INIT {
    __TOKEN_PROVIDER = loadTokenPool() 
  }

  return __TOKEN_PROVIDER
}

func loadTokenPool() *TokenPool {
  __TOKEN_INIT = false

  tokenPool := &TokenPool{ 
    pool: pool{FilePath: lib.Sys.TokenFile},
  }
  tokenPool.self = tokenPool
  tokenPool.Build()

  return tokenPool
}

var (
  __TOKEN_PROVIDER *TokenPool
  __TOKEN_INIT bool = true
)
