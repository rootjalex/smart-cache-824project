package task 

import (
	"strconv"
	// "testing"

	"../config"
	"../datastore"
	// "../utils"

	// "log"
	"fmt"
)

func MakeDatastore(numDatastoreFiles int) (*datastore.DataStore, map[string]config.DataType, []string, []config.DataType) {
	ds := datastore.MakeDataStore()
	files := make(map[string]config.DataType)
	fileNames, fileContents := []string{}, []config.DataType{}

	// add files to the datastore
	for i := 0; i < numDatastoreFiles; i++ {
		fileName := "MNIST_" + strconv.Itoa(i+1) + ".png"
		fileContent := (config.DataType)("image_" + strconv.Itoa(i+1))
		ds.Make(fileName, fileContent)

		files[fileName] = fileContent
		fileNames = append(fileNames, fileName)
		fileContents = append(fileContents, fileContent)
	}

	return ds, files, fileNames, fileContents
}

func PrintFailure(failed bool) {
	if failed {
		fmt.Printf("\t... FAILED\n")
	} else {
		fmt.Printf("\t... PASSED\n")
	}
}