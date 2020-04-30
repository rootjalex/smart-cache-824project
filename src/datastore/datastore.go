package datastore

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

type DataType string

type DataStore struct {
    mu     sync.Mutex
    data   map[string]DataType
    n      int
}

func MakeDataStore(data map[string]DataType) *DataStore {
    d := &DataStore{}
    d.data = data
    d.n = len(data)
    return d
}

func (d *DataStore) Size() int {
    return d.n
}

func (d *DataStore) Get(filename string) DataType {
    // TODO: add time.Sleep for approx time of fetching from underlying datastore
    return d.data[filename]
}

func (d *DataStore) Make(filename string) {
    // TODO: uncomment if DataType is *os.File
    // f, _ := os.Create(filename) // ignore error if already exists
	// f.Close()
    // d.data[filename], _ = os.Open(filename)

    // TODO: comment out if DataType is *os.File
    d.data[filename] = "good"
}