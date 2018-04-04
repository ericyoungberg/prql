package lib

import (
  "fmt"
  "strings"
  "syscall"
  "crypto/md5"
  "encoding/hex"

  "golang.org/x/crypto/ssh/terminal"
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
