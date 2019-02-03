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
  testStrings = map[string]string{
    "a key": "ec4f6067|32e9ea4aae54cf8d195f6b34edc3f4cc30994d164d005b0150cfe199a0aaddbfe0",
    "testkey2": "a2e27a3d|e592bbc3d8e1efac90223630f29baedf327b4a198abd738b1e702fadc58630872288df85", 
    "Streaming a bad dystopian movie right now": "7bbf55de|cba5f9822499f2eac97f2cc957da2b3ca56c0d8fbae8908c1e44dc6b89db4f9d6867b5fb179780f423bb2ac3e89af76f66f39ae94caef6cbc7ed20d65787c251a642850075", 
    "areallylongkeythatwillreallypushthelimitsbeyondrealworldusecases": "aa571bac|470178ce24bfc4584ee54f9934d5ea0c17757f798a03be535f53b4a88209145e2dd29a2d218f9c2fb6bdb5db1baee94f0cfaf937ecd3d960ba9211ca0c9b78599a6b7a20f219d88577ef6bf29f21bd38ca4b2ce0a3095e201d4359a7",
  }
)


func TestInsecureCryptography(t *testing.T) {
  for data, _ := range testStrings {
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
  for key, encrypted := range testStrings {
    decrypted, err := InsecureDecryptString(encrypted)

    if err != nil {
      t.Error(err) 
    } else if decrypted != key {
      t.Errorf("expected decrypted data to be %s, was %s [value=%s]", key, decrypted, encrypted) 
    }
  } 
}

func TestInsecureEncryptString(t *testing.T) {
  encryptionFmt := fmt.Sprintf("^\\w{%d,}\\%s\\w+$", minSaltLength, saltDelimiter)
  encryptionRegex, _ := regexp.Compile(encryptionFmt)

  for data, _ := range testStrings {
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

  for key, _ := range testStrings {
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
