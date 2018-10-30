package lib

import (
  "io"
  "fmt"
  "strings"
  "syscall"
  "net/http"
  "crypto/aes"
  "crypto/md5"
  "crypto/rand"
  "crypto/cipher"
  "encoding/hex"

  "golang.org/x/crypto/ssh/terminal"
  log "github.com/sirupsen/logrus"
)


func CreateHash(key string) string {
  hasher := md5.New()
  hasher.Write([]byte(key))
  return hex.EncodeToString(hasher.Sum(nil))
}


func Encrypt(data []byte, passphrase string) []byte {
  block, _ := aes.NewCipher([]byte(CreateHash(passphrase)))
  gcm, err := cipher.NewGCM(block)
  if err != nil {
    panic(err.Error()) 
  }

  nonce := make([]byte, gcm.NonceSize())

  if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
    panic(err.Error())
  }

  cipherText := gcm.Seal(nonce, nonce, data, nil)

  return cipherText
}


func Decrypt(data []byte, passphrase string) []byte {
  key := []byte(CreateHash(passphrase))
  block, err := aes.NewCipher(key)

  if err != nil {
    panic(err.Error()) 
  }

  gcm, err := cipher.NewGCM(block)
  if err != nil {
    panic(err.Error()) 
  }

  nonceSize := gcm.NonceSize()
  nonce, cipherText := data[:nonceSize], data[nonceSize:]
  plainText, err := gcm.Open(nil, nonce, cipherText, nil)
  if err != nil {
    panic(err.Error()) 
  }

  return plainText
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
      http.Error(w, errorMessage, http.StatusBadRequest)
    }
  }
}
