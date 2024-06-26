package gohtslib

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

func ReadSiteCandidates(filename string) (map[string][]string, error) {
	fmt.Printf("Debug : marking ReadSiteCandidates : \"%s\"\n", filename)
	if strings.HasSuffix(strings.ToLower(filename), ".csv") {
		df := DataFrameFromCSV[string](filename)
		if !SetEq(df.Heads(), append(PURE_SITES, "Element")) {
			return nil, fmt.Errorf("file does not contain the right headers, must be \"Element\", \"A\", \"B\" and \"C\", not %v != %v", df.Heads(), append(PURE_SITES, "Element"))
		}
		result := make(map[string][]string, len(PURE_SITES))
		for _, site := range PURE_SITES {
			result[site] = []string{}
		}
		for _, row := range df.IterRows() {
			for _, site := range PURE_SITES {
				if IsIn(row[site], []string{"yes", "y"}) {
					result[site] = append(result[site], row["Element"])
				} else if !IsIn(row[site], []string{"no", "n"}) {
					return nil, fmt.Errorf("Value should be \"yes\", \"y\", \"no\" or \"n\", not \""+row[site]+"\"")
				}
			}
		}
		return result, nil
	} else if strings.HasSuffix(strings.ToLower(filename), ".json") {
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
			if !SetEq(Keys(data), append(SITES, "Element")) {
				return nil, fmt.Errorf("file does not contain the right headers, must be \"Element\", \"A\", \"B\" and \"C\"")
			}
			return data, nil
		}
		var list []string
		if err := json.Unmarshal(byteValue, &list); err == nil {
			// Convert the list to a map with a default key
			return map[string][]string{"default": list}, nil
		}
		return nil, fmt.Errorf("failed to unmarshal JSON")
	} else {
		panic(fmt.Sprintf("Unknown extension for file, must be json or csv. File : %s", filename))
		return nil, nil
	}
}
