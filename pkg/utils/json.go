package utils

import (
	"encoding/json"
	"os"
)

func JsonUnmarshalStrict(file string, obj interface{}) error {
	// open the file
	f, err := os.Open(file)
	if err != nil {
		return err
	}
	// create a decoder for the files content
	decoder := json.NewDecoder(f)
	// enable strict parsing
	decoder.DisallowUnknownFields()
	// unmarshall
	err = decoder.Decode(obj)
	if err != nil {
		return err
	}
	return nil
}
