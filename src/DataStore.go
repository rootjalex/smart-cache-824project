package cache

import (
    "sync"
)

type DataStore struct {
    mu     sync.Mutex
    data   []interface{}
    n      int
}

func MakeDataStore(data []interface{}) *DataStore {
    d := &DataStore{}
    d.data = data
    d.n = len(data)
    return d
}

func (d *DataStore) Size() int {
    return d.n
}

