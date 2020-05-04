package datastore

import (
    "sync"
    "../config"
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
    data   map[string]config.DataType
    n      int
}

func (d *DataStore) GetFileNames() []string {
    filenames := make([]string, len(d.data))
    i := 0
    for f, _ := range d.data {
        filenames[i] = f
        i++
    }
    return filenames
}

func MakeDataStore() *DataStore {
    d := &DataStore{}
    d.data = make(map[string]config.DataType)
    d.n = 0
    return d
}

func (d *DataStore) Size() int {
    return d.n
}

func (d *DataStore) Get(filename string) (config.DataType, bool) {
    // TODO: add time.Sleep for approx time of fetching from underlying datastore
    data, ok := d.data[filename]
    return data, ok
}

func (d *DataStore) Make(filename string) {
    // TODO: uncomment if DataType is *os.File
    // f, _ := os.Create(filename) // ignore error if already exists
	// f.Close()
    // d.data[filename], _ = os.Open(filename)

    // TODO: comment out if DataType is *os.File
    d.data[filename] = "good"

    d.n = len(d.data)
}

func (d *DataStore) Copy() *DataStore {
    c := &DataStore{}
    c.data = make(map[string]config.DataType)
    c.n = 0
    return c
}
