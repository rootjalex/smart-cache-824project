package cache

import (
    "sync"
)
/********************************************************
DataStore API

Make(data string)
 - intializes a datastore storing the inut data

Size()
 - returns the size (number of files) in the datastore

Get(file string)
 - returns the datastore in the file for the corresponding key
********************************************************/

type DataStore struct {
    mu     sync.Mutex
    data   map[string]string
    n      int
}

func MakeDataStore(data map[string]string) *DataStore {
    d := &DataStore{}
    d.data = data
    d.n = len(data)
    return d
}

func (d *DataStore) Size() int {
    return d.n
}

func (d *DataStore) Get(f string) string {
    return d.data[f]
}

