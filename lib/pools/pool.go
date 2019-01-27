package pools

import (
  "strings"
  "io/ioutil"

  log "github.com/sirupsen/logrus"
)

const (
  entryDelimiter string = ":"
)

type Pool interface {
  build()
}

type pool struct {
  child Pool
  FilePath string

  records [][]string
}

func (p *pool) Save() error {
  lines := make([]string, len(p.records))

  for i, entry := range p.records {
    lines[i] = strings.Join(entry, entryDelimiter)
  } 

  data := []byte(strings.Join(lines, "\n"))

  return ioutil.WriteFile(p.FilePath, data, 0600)
}

func (p *pool) Remove(keys []string) {
  p.records = removeByColumn(keys, p.records, 0)
  p.child.build()
}

func (p *pool) AppendRecord(record []string) {
  p.records = append(p.records, record)
  p.child.build()
}

func (p *pool) Build() {
  p.load()
  p.child.build()
}

func (p *pool) load() {
  entries := ParseEntryFile(p.FilePath)
  p.records = make([][]string, len(entries))
  p.records = entries
}

func (p *pool) build() {
  log.Fatal("pool.build() must be overriden")
}

