package lib

import (
  "fmt"
  "strings"
  "syscall"
  "net/http"
  "crypto/md5"
  "encoding/hex"

  "golang.org/x/crypto/ssh/terminal"
  log "github.com/sirupsen/logrus"
)


func CreateHash(seed string) string {
  hasher := md5.New()
  hasher.Write([]byte(seed))
  return hex.EncodeToString(hasher.Sum(nil))
}


func GetPassword(user string) (string, error) {
  if user != "" {
    fmt.Print(user + "'s password: ")
  } else {
    fmt.Print("Password: ") 
  }

  defer fmt.Print("\n")

  input, err := terminal.ReadPassword(int(syscall.Stdin))
  if err != nil {
    return "", err 
  }


  return strings.TrimSpace(string(input)), nil
}

func SecretExec(fn func()) func(http.ResponseWriter, *http.Request) {
  return func(w http.ResponseWriter, r *http.Request) {
    config, err := GetConfig()
    if err != nil {
      log.Fatal(err) 
    }

    clientSecret := r.Header.Get(config.Headers.Secret)
    serverSecret := config.Secret

    if clientSecret == serverSecret {
      fn()
    } else {
      errorMessage := "command is only available to local prql"
      ipLogger := NewIPLogger(r) 
      ipLogger.Error(errorMessage)
      http.Error(w, errorMessage, http.StatusBadRequest)
    }
  }
}
