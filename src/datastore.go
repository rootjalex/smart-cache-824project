package cache

import (
    "sync"
    "os"
)
/************************************************
DataStore API

Make(data)
 - intializes a datastore storing the inut data

Size()
 - returns the size (number of files) in the datastore
************************************************/

type DataStore struct {
    mu     sync.Mutex
    data   []*os.File
    n      int
}

func MakeDataStore(data []*os.File) *DataStore {
    d := &DataStore{}
    d.data = data
    d.n = len(data)
    return d
}

func (d *DataStore) Size() int {
    return d.n
}

