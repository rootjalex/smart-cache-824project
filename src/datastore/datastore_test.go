package datastore

import (
    "fmt"
    "testing"
)

func TestDatastoreCopy(t *testing.T) {
	fmt.Println("TestDatastoreCopy ...")
	// make datastore
    d := MakeDataStore()

    d.Make("1", "hi")
    d.Make("2", "bye")
    c := d.Copy()

    first, _:= c.Get("1")
    second, _ := c.Get("2")
    if first != "hi" || second != "bye" {
	    t.Errorf("FAILED copied datastore: %v, original: %v", c, d)

    }

}

