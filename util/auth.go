package util

import (
  "crypto/md5"
  "encoding/hex"
)


func CreateHash(seed string) string {
  hasher := md5.New()
  hasher.Write([]byte(seed))
  return hex.EncodeToString(hasher.Sum(nil))
}


func GetPassword() string {
  return ""
}
