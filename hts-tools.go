//@hts
//do not delete previous comment

package gohtslib

import (
	"fmt"
	"strings"
)

func ExtractMaterial(line string) map[string]any {
	isNumeric := func(r rune) bool {
		return strings.ToLower(string(r)) == strings.ToUpper(string(r))
	}
	isUpper := func(r rune) bool {
		return strings.ToUpper(string(r)) == string(r) && !isNumeric(r)
	}
	isLower := func(r rune) bool {
		return strings.ToLower(string(r)) == string(r) && !isNumeric(r)
	}
	result := map[string]any{}
	runes := []rune(line)
	i := 0
	for s := range PURE_SITES {
		if !isUpper(runes[i]) {
			panic(fmt.Sprintf("Site %s does not start with upper case : %s", PURE_SITES[s], line))
		}
		el_x := string(runes[i:i+1])
		i += 1
		for isLower(runes[i]) {
			el_x += string(runes[i])
			i += 1
		}
		if isNumeric(runes[i]) {
			cx := ""
			for isNumeric(runes[i]) {
				cx += string(runes[i])
				i += 1
			}
			if !isUpper(runes[i]) {
				panic(fmt.Sprintf("Site %s does not start with upper case : %s", DOPING_SITES[s], line))
			}
			el_xs := string(runes[i:i+1])
			i += 1
			for isLower(runes[i]) {
				el_xs += string(runes[i])
				i += 1
			}
			if !isNumeric(runes[i]) {
				panic(fmt.Sprintf("Site %s has no concentration : %s", DOPING_SITES[s], line))
			}
			cxs := ""
			for i<len(runes) && isNumeric(runes[i]) {
				cxs += string(runes[i])
				i += 1
			}
			result[PURE_SITES[s]] = el_x
			result[DOPING_SITES[s]] = el_xs
			result[CONCENTRATIONS[s]] = stringToFloat64(cxs)
		}
	}
	return result
}

func VectorMaterial(material map[string]any, T float64, table map[string]map[string]float64, keys []string) []float64 {
	if keys == nil {
		keys = Apply(
			func(elementAndSite []string) string {
				return elementAndSite[0] + elementAndSite[1]
			},
			Product(Keys(table), SITES),
		)
	}
	keys = Filter(func(key string)bool{return !(key=="T"||key=="CAs"||key=="CBs"||key=="CCs")}, keys)
	result := make([]float64, 4+len(keys))
	result[0] = T
	result[1] = material["CAs"].(float64)
	result[2] = material["CBs"].(float64)
	result[3] = material["CCs"].(float64)
	for k, key := range keys {
		result[4+k] = table[CropSite(key)][material[GetSite(key)].(string)]
	}
	return result
}

func PredictMaterialLine(input string, T float64, model Model, table map[string]map[string]float64) float64 {
	material := ExtractMaterial(input)
	vector := VectorMaterial(material, T, table, model.keys)
	fmt.Printf("%v\n", vector)
	result := model.predictRows([][]float64{vector})[0]
	return result
}

func LoadMatrixTable(pfs string, elements []string, keys []string) [][]float64 {
	df := DataFrameFromCSV[string](pfs)
	if keys == nil {
		keys = df.Heads()
	}
	for name := range df.columns {
		found := false
		for _, site := range SITES {
			for _, key := range keys {
				if key == name + site {
					found = true
					break
				}
			}
			if found {
				break
			}
		}
		if !found && name != "Element" {
			delete(df.columns, name)
		}
	}
	df.heads = df.Heads()
	fmt.Printf("Table : \n%s\n", df)
	table := make([][]float64, len(keys))
	for k, key := range keys {
		table[k] = make([]float64, len(elements))
		for e, element := range elements {
			if element == "Du" {
				continue
			}
			table[k][e] = stringToFloat64(df.columns[CropSite(key)][IndexPanic(element, df.columns["Element"])])
		}
	}
	return table
}

func LoadMapTable(pfs string, keys []string) map[string]map[string]float64 {
	df := DataFrameFromCSV[string](pfs)
	if keys == nil {
		keys = df.Heads()
	}
	table := make(map[string]map[string]float64, len(keys))
	for _, key := range keys {
		table[key] = make(map[string]float64, df.Length())
	}
	for _, row := range df.IterRows() {
		for _, key := range keys {
			table[key][row["Element"]] = stringToFloat64(row[CropSite(key)])
		}
	}
	return table
}

func IterateOverClasses(dataframe *DataFrame[int], table [][] float64, keys []string, concentrations []float64, temperature float64) ([]string, [][]float64) {
	binary := []string{"000", "001", "010", "011", "100", "101", "110", "111"}
	nbClasses := dataframe.Length()
	size := 0
	println("Building NbPermutations...")
	dataframe.AddColumn("NbPermutations", Apply(
		func(row map[string]int) int {
			result := 1
			for _, x := range PURE_SITES {
				if row[x] != row[x+"s"] {
					result *= len(concentrations)
				}
			}
			size += result
				return result
		},
		dataframe.IterRows(),
	))
	fmt.Printf("\x1b[93m`size` -> %v\x1b[39m\n", size)
	println("Building DopingProfile...")
	dataframe.AddColumn("DopingProfile", Apply(
		func(row map[string]int) int {
			key := ""
			for _, x := range PURE_SITES {
				if row[x] == row[x+"s"] {
					key += "0"
				} else {
					key += "1"
				}
			}
			return int(Index(key, binary))
		},
		dataframe.IterRows(),
	))
	println("Building keySites...")
	keySites := Apply(
		func(key string) string {
			for _, site := range SITES {
				if strings.HasSuffix(key, site) {
					return site
				}
			}
			panic("key without site : "+key)
			return ""
		},
		keys,
	)
	println("Building permutations...")
	permutations := Apply(
		func(i int) [][]float64 {
			return Transpose(Product(Apply(
				func(r rune) []float64 {
					if r == '0' {
						return []float64{0}
					} else {
						return concentrations
					}
				},
				[]rune(binary[i]),
			)...))
		},
		Range(8),
	)
	labels := make([]string, size)
	p := 0
	for i:=0; i<nbClasses; i++ {
		fmt.Printf("\rBuilding class label %d / %d", i, nbClasses)
		n := dataframe.columns["NbPermutations"][i]
		name := ""
		for _, site := range SITES {
			name += fmt.Sprintf("%d:", dataframe.columns[site][i])
		}
		for j:=0; j<n; j++ {
			labels[p+j] = fmt.Sprintf("%s%d", name, j)
		}
		p += n
	}
	println("")
	input := make([][]float64, 4+len(keys))
	input[0] = Repeat(temperature, size)
	// fmt.Printf("keys     = %v\nkeySites = %v\n", keys, keySites)
	for i:=0; i<nbClasses; i++ {
		for p, permutationColumn := range permutations[dataframe.columns["DopingProfile"][i]] {
			if len(permutationColumn) != dataframe.columns["NbPermutations"][i] {
				panic(fmt.Sprintf("Adding line #%d : %v; len(permutationColumn)=%v; dataframe.columns[\"NbPermutations\"][i]=%v", i, dataframe.IterRows()[i], len(permutationColumn), dataframe.columns["NbPermutations"][i]))
			}
			input[1+p] = append(input[1+p], permutationColumn...)
		}
		for k := range keys {
			input[4+k] = append(
				input[4+k], 
				Repeat(
					table[k][dataframe.columns[keySites[k]] [i]],
					dataframe.columns["NbPermutations"][i],
				)...,
			)
		}
		fmt.Printf("\rClass : %d / %d", i+1, nbClasses)
	}
	println("")
	return labels, input
}
