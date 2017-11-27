package auth

import (
  "fmt"
  "crypto/hash"
)

func GenerateToken() string {
  return hash.MD5.New()
}
