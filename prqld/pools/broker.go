package pool

import (
  "github.com/prql/prql/lib"
)

type Pool map[string]lib.Entry

type PoolBroker struct {
  pools map[string]Pool
}

func (pb *PoolBroker) GetPool(pool string) Pool {
  return pb.pools[pool]
}

func (pb *PoolBroker) Get(pool string, key string) lib.Entry {
  return pb.GetPool(pool)[key]
}
 
func (pb *PoolBroker) Add(pool string, key string, entry lib.Entry) {
  pb.pools[pool][key] = entry
}


func GetPoolBroker() *PoolBroker {
  if __PROVIDER.pools == nil {
    __PROVIDER = PoolBroker{}
  }

  return &__PROVIDER
}

var (
  __PROVIDER PoolBroker
)
