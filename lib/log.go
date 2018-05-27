package lib

import (
  "net/http"

  log "github.com/sirupsen/logrus"
)


func NewIPLogger(r *http.Request) *log.Entry {
  return log.WithFields(log.Fields{"IP": r.RemoteAddr})
}
