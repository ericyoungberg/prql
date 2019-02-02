package lib

import (
  "fmt"
  "regexp"
  "strings"
  "testing"
)

const (
  minSaltLength = 4
)

var (
  testStrings = [4]string{
    "a key", 
    "testkey2", 
    "Streaming a bad dystopian movie right now", 
    "areallylongkeythatwillreallypushthelimitsbeyondrealworldusecases",
  }
)


func TestInsecureCryptography(t *testing.T) {
  for _, data := range testStrings {
    encrypted := InsecureEncryptString(data)
    decrypted, err := InsecureDecryptString(encrypted)

    if err != nil  {
      t.Error(err)    
    }

    if data != decrypted {
      t.Errorf("insecure encryption failure: [value=%s] [encrypted=%s] [decrypted=%s]", data, encrypted, decrypted)
    }
  }
}

func TestInsecureDecryptString(t *testing.T) {
  for _, data := range testStrings {
    
  } 
}

func TestInsecureEncryptString(t *testing.T) {
  encryptionFmt := fmt.Sprintf("^\\w{%d,}\\%s\\w+$", minSaltLength, saltDelimiter)
  encryptionRegex, _ := regexp.Compile(encryptionFmt)

  for _, data := range testStrings {
    encryptedData := InsecureEncryptString(data)

    if !encryptionRegex.MatchString(encryptedData) {
      parts := strings.Split(encryptedData, saltDelimiter)

      if len(parts[0]) < minSaltLength {
        t.Errorf("expected salt %s to have minimum length of %d", parts[0], minSaltLength)
      } else {
        t.Errorf("encrypted data output improperly formatted. Should be <salt>%s<data>. [value=%s]", saltDelimiter, encryptedData) 
      }
    }
  }
}

func TestCreateHash(t *testing.T) {
  hashes := make(map[string]string)

  for _, key := range testStrings {
    hash := CreateHash(key)
    
    if hash == key {
      t.Errorf("key equals produced hash: %s == %s", hash, key) 
    } else if found, ok := hashes[hash]; ok {
      t.Errorf("duplicate hash %s from keys %s and %s", hash, key, found) 
    } else {
      hashes[hash] = key
    }
  }
}

func TestCreateSalt(t *testing.T) {
  saltSizes := []int{1, 4, 25, 0}

  for _, size := range saltSizes {
    if salt, err := createSalt(size); err != nil {
      t.Errorf("couldn't create salt of size %d", size) 
    } else {
      if len(salt) != size {
        t.Errorf("expected salt to have size %d, got: %d [value=%s]", size, len(salt), salt) 
      }
    }
  }
}
