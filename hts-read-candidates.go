package gohtslib

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

func readSiteCandidates(filename string) (map[string][]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %v", err)
	}	
	defer file.Close()
	byteValue, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %v", err)
	}
	var data map[string][]string
	if err := json.Unmarshal(byteValue, &data); err == nil {
		return data, nil
	}
	var list []string
	if err := json.Unmarshal(byteValue, &list); err == nil {
		// Convert the list to a map with a default key
		return map[string][]string{"default": list}, nil
	}
	return nil, fmt.Errorf("failed to unmarshal JSON")
}
