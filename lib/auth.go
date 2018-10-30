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


const (
  saltDelimiter string =  "|"
)


func CreateHash(key string) string {
  hasher := md5.New()
  hasher.Write([]byte(key))
  return hex.EncodeToString(hasher.Sum(nil))
}


func InsecureEncryptString(data string) string {
  salt, err := createSalt(4)
  if err != nil {
    log.Fatal("could not create salt for string encryption")  
  }
  cleanSalt := hex.EncodeToString(salt)
  encrypted := encrypt([]byte(data), cleanSalt)

  return strings.Join([]string{cleanSalt, hex.EncodeToString(encrypted)}, saltDelimiter)
}

func InsecureDecryptString(data string) (string, error) {
  pieces := strings.Split(data, saltDelimiter)
  salt, encryptedHex := pieces[0], pieces[1]

  encrypted, err := hex.DecodeString(encryptedHex)
  if err != nil {
    return "", err 
  }
  
  return string(decrypt(encrypted, salt)), nil
}


func encrypt(data []byte, passphrase string) []byte {
  block, _ := aes.NewCipher([]byte(CreateHash(passphrase)))
  gcm, err := cipher.NewGCM(block)
  if err != nil {
    panic(err.Error()) 
  }

  nonce, err := createSalt(gcm.NonceSize())
  if err != nil {
    log.Fatal("could not create nonce for encryption")
  }

  cipherText := gcm.Seal(nonce, nonce, data, nil)

  return cipherText
}


func decrypt(data []byte, passphrase string) []byte {
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


func GetPasswordFromTerminal(user string) (string, error) {
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

func createSalt(size int) ([]byte, error) {
   salt := make([]byte, size)

   if _, err := io.ReadFull(rand.Reader, salt); err != nil {
    return nil, err
  }

  return salt, nil
}

