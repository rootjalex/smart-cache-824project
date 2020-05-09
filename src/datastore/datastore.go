package datastore

import (
    "sync"
    "time"
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
    mu      sync.Mutex
    data    map[string]config.DataType
    n       int
    calls   int64
}

func (d *DataStore) CountCalls() int64 {
    d.mu.Lock()
    defer d.mu.Unlock()
    return d.calls
}

func (d *DataStore) GetFileNames() []string {
    d.mu.Lock()
    defer d.mu.Unlock()
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
    d.calls = 0
    return d
}

func (d *DataStore) Size() int {
    d.mu.Lock()
    defer d.mu.Unlock()
    return d.n
}

func (d *DataStore) Get(filename string) (config.DataType, bool) {
    time.Sleep(config.DATA_FETCH_TIME)
    d.mu.Lock()
    defer d.mu.Unlock()
    data, ok := d.data[filename]
    d.calls++
    // approx time of fetching from underlying datastore
    return data, ok
}

func (d *DataStore) Make(filename string, content config.DataType) {
    d.mu.Lock()
    defer d.mu.Unlock()

    // TODO: uncomment if DataType is *os.File
    // f, _ := os.Create(filename) // ignore error if already exists
	// f.Close()
    // d.data[filename], _ = os.Open(filename)

    // TODO: comment out if DataType is *os.File
    d.data[filename] = content
    d.n = len(d.data)
}

func (d *DataStore) Copy() *DataStore {
    d.mu.Lock()
    defer d.mu.Unlock()
    c := &DataStore{}
    c.data = make(map[string]config.DataType)
    for filename, content := range d.data {
        c.data[filename] = content
    }
    c.n = len(c.data)
    c.calls = 0
    return c
}
