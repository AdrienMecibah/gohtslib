//@hts
//do not delete previous comment

package gohtslib

import (
	"strings"
	"fmt"
)

func ReadMapTable(filepath string) map[string]map[string]float64 {
	table := DataFrameFromCSV[string](filepath)
	mapTable := Map(
		func(key string) (string, map[string]float64) {
			return key, Map(
				func(row map[string]string) (string, float64) {
					return row["Element"], stringToFloat64(row[key])
				},
				table.IterRows(),
			)
		},
		Filter(func(h string)bool{return h!="Element"}, table.heads),
	)
	return mapTable
}

func WriteClassName(row map[string]string) string {
	result := ""
	for _, x := range SITES {
		if row[x] == row[x+"s"] || row[x+"s"] == "Du" {
			result += row[x]
		} else {
			result += fmt.Sprintf("(%s%s)", row[x], row[x+"s"])
		}
	}
	return result
}

func CropSite(name string) string {
	type Runes []rune
	for _, site := range SITES {
		if strings.HasSuffix(name, site) {
			return string(Runes(name)[:len(Runes(name))-len(Runes(site))])
		}
	}
	return name
}

func GetSite(name string) string {
	for _, site := range SITES {
		if strings.HasSuffix(name, site) {
			return site
		}
	}
	panic("name \""+name+"\" does not end with any sites")
	return ""
}

func HasSite(name string) bool {
	for _, site := range SITES {
		if strings.HasSuffix(name, site) {
			return true
		}
	}
	return false
}


func UnLabel(label string, elements []string, concentrations []float64) string {
	result := ""
	labelElements := StrSplit(label, ":")
	iterations := make([][]float64, 3)
	for i := range PURE_SITES {
		if labelElements[2*i] == labelElements[2*i+1] {
			iterations[i] = []float64{0}
		} else {
			iterations[i] = concentrations
		}
	}
	permutations := Product(iterations...)
	concentrationValues := permutations[stringToInt(labelElements[len(labelElements)-1])]
	for i := range PURE_SITES {
		if labelElements[2*i] == labelElements[2*i+1] {
			result += elements[stringToInt(labelElements[2*i])]
		} else {
			result += fmt.Sprintf("%s%v%s%v", elements[stringToInt(labelElements[2*i])], 1-concentrationValues[i], elements[stringToInt(labelElements[2*i+1])], concentrationValues[i])
		}
	}
	return result
}


func MultipleUnLabel(labels []string, elements []string, concentrations []float64) []string {
	result := make([]string, len(labels))
	binary := []string{"000", "001", "010", "011", "100", "101", "110", "111"}
	permutations := Map(
		func(i int) (string, [][]float64) {
			return binary[i], Product(Apply(
				func(r rune) []float64 {
					if r == '0' {
						return []float64{0}
					} else {
						return concentrations
					}
				},
				[]rune(binary[i]),
			)...)
		},
		Range(8),
	)
	for i, label := range labels {
		fmt.Printf("\rBuilding material name %d / %d", i+1, len(labels))
		labelElements := StrSplit(label, ":")
		bin := ""
		for i := range PURE_SITES {
			if labelElements[2*i] == labelElements[2*i+1] {
				bin += "0"
			} else {
				bin += "1"
			}
		}
		concentrationValues := permutations[bin][stringToInt(labelElements[len(labelElements)-1])]
		material := ""
		for i := range PURE_SITES {
			if labelElements[2*i] == labelElements[2*i+1] {
				material += elements[stringToInt(labelElements[2*i])]
			} else {
				material += fmt.Sprintf("%s%v%s%v", elements[stringToInt(labelElements[2*i])], 1-concentrationValues[i], elements[stringToInt(labelElements[2*i+1])], concentrationValues[i])
			}
		}
		result[i] = material
	}
	println("")
	return result
}
