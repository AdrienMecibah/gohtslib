//@hts
//do not delete previous comment

package gohtslib

import (
	"fmt"
)

var _ = Script("list-models", func(){
	for name := range models {
		println(name)
	}
})

var _ = Script("test-model", func(){
	argv := ParseArgv[struct{
		model string
		table string
		input string
		output string
	}]()
	if !argv.present["input"] {
		panic(fmt.Sprintf("Flag -input is mandatory"))
	}
	if !argv.present["table"] {
		panic(fmt.Sprintf("Flag -table is mandatory"))
	}
	if !argv.present["model"] {
		panic(fmt.Sprintf("Flag -model is mandatory : must be one of %v", models))
	}
	model, found := models[argv.flags.model]
	if !found {
		panic(fmt.Sprintf("Unknown model \"%s\". Must be one of %v", argv.flags.model, Repr(Keys(models))))
	}
	dataset := DataFrameFromCSV[string](argv.flags.input)
	dataset = DataFrameFromRows(Apply(
		func(row map[string]string) map[string]string {
			for i := range PURE_SITES {
				if row[DOPING_SITES[i]] == "Du" {
					row[DOPING_SITES[i]] = row[PURE_SITES[i]]
				}
			}
			return row
		},
		dataset.IterRows(),
	))
	table := DataFrameFromCSV[string](argv.flags.table)
	table = table.SelectRows(func(row map[string]string)bool{
		for _, site := range SITES {
			for _, datasetRow := range dataset.IterRows() {
				if datasetRow[site] == row["Element"] {
					return true
				}
			}
		}
		return false
	})
	table = table.SelectColumns(func(name string, column[]string)bool{
		if name == "Element" {
			return true
		}
		for _, site := range SITES {
			if IsIn(name+site, model.Keys) {
				return true
			}
		}
		return false
	})
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
	fmt.Printf("table.heads = %v\n", table.heads)
	fmt.Printf("mapTable = %v\n", mapTable)
	input := Apply(
		func(row map[string]string) []float64 {
			newRow := map[string]any{}
			for _, site := range SITES {
				newRow[site] = any(row[site])
			}
			for _, cxs := range CONCENTRATIONS {
				newRow[cxs] = any(stringToFloat64(row[cxs]))
			}
			fmt.Printf("\x1b[38;5;226m%v\x1b[39m\n", newRow)
			return VectorMaterial(newRow, stringToFloat64(row["T"])/float64(1200), mapTable, model.Keys)
		}, 
		dataset.IterRows(),
	)
	pred := model.PredictRows(input)
	target := Apply(stringToFloat64, dataset.columns["ZT"])
	dataset.AddColumn("Prediction", Apply(float64ToString, pred))
	fields := []string{"A", "As", "B", "Bs", "C", "Cs", "T", "ZT", "Prediction"}
	dataset = dataset.SelectColumns(func(name string, column []string)bool{return IsIn(name, fields)})
	if argv.present["output"] {
		dataset.Save(argv.flags.output)
	}
	fmt.Printf("%v\nRMSE :%v\n", dataset, Rmse(target, pred))
})

var _ = Script("predict", func(){
	argv := ParseArgv[struct{
		model string
		table string
		material string
		input string
		T float64
	}]()
	if !argv.present["material"] && !argv.present["input"] {
		panic(fmt.Sprintf("Using either -material or -input flag is mandatory"))
	}
	if !argv.present["T"] {
		panic(fmt.Sprintf("Flag -T is mandatory", models))
	}
	if !argv.present["table"] {
		panic(fmt.Sprintf("Flag -table is mandatory"))
	}
	if !argv.present["model"] {
		panic(fmt.Sprintf("Flag -model is mandatory : must be one of %v", models))
	}
	model, found := models[argv.flags.model]
	if !found {
		panic(fmt.Sprintf("Unknown model \"%s\". Must be one of %v", argv.flags.model, Repr(Keys(models))))
	}
	material := ExtractMaterial(argv.flags.material)
	// keys := Filter(func(key string)bool{return !(key=="T"||key=="CAs"||key=="CBs"||key=="CCs")}, model.Keys)
	table := DataFrameFromCSV[string](argv.flags.table)
	table = table.SelectRows(func(row map[string]string)bool{
		for _, site := range SITES {
			if element, ok := material[site].(string); ok && element == row["Element"] {
				return true
			}
		}
		return false
	})
	table = table.SelectColumns(func(name string, column[]string)bool{
		if name == "Element" {
			return true
		}
		for _, site := range SITES {
			if IsIn(name+site, model.Keys) {
				return true
			}
		}
		return false
	})
	table.heads = append([]string{"Element"}, WithoutDuplicates(append(table.heads, "Element"))...)
	fmt.Printf("%s\n", table)
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
	fmt.Printf("%v\n", PredictMaterialLine(argv.flags.material, argv.flags.T, model, mapTable))
})

var _ = Script("search", func() {

	Step("Global process", func(){

		argv := ParseArgv[struct{
			classBuilder string
			table string
			builderArgs string
			model string
			outputSize int
			output string
			siteCandidates string
		}]()

		if !argv.present["classBuilder"] {
			panic(fmt.Sprintf("Flag -class-builder is mandatory : must be one of %v", classesBuilders))
		}
		if !argv.present["table"] {
			panic(fmt.Sprintf("Flag -table is mandatory"))
		}
		if !argv.present["model"] {
			panic(fmt.Sprintf("Flag -model is mandatory : must be one of %v", models))
		}
		if !argv.present["outputSize"] {
			panic(fmt.Sprintf("Flag -output-size is mandatory"))
		}
		if !argv.present["output"] {
			panic(fmt.Sprintf("Flag -output is mandatory"))
		}
	
		if argv.present["siteCandidates"] {
			candidates, err := ReadSiteCandidates(argv.flags.siteCandidates)
			if err != nil {
				panic(err)
			}
			if !SetEq(Keys(candidates), PURE_SITES) {
				panic(fmt.Sprintf("Keys in file %s must be %v, not %v", argv.flags.siteCandidates, PURE_SITES, Keys(candidates)))
			}
			CANDIDATES = []string{}
			for _, site := range PURE_SITES {
				for _, x := range candidates[site] {
					if !IsIn(x, CANDIDATES) {
						CANDIDATES = append(CANDIDATES, x)
					}
				}
			}
			A_CANDIDATES = candidates["A"]
			B_CANDIDATES = candidates["B"]
			C_CANDIDATES = candidates["C"]
		}

		model, found := models[argv.flags.model]
		if !found {
			panic(fmt.Sprintf("Unknown model \"%s\". Must be one of %v", argv.flags.model, Repr(Keys(models))))
		}
		fmt.Printf("model \"%s\" : %v\n", argv.flags.model, model)

		concentrations := []float64{0.01, 0.02, 0.03, 0.04, 0.05, 0.06, 0.07, 0.08, 0.09, 0.10, 0.11, 0.12, 0.13, 0.14, 0.15, 0.16, 0.17, 0.18, 0.19, 0.20, 0.30, 0.40, 0.50}
		var elements []string
		var dataset  *DataFrame[int]
		var keys     []string 
		var table    [][]float64	
		var labels    []string
		var input     [][]float64
		var preds     []float64
		var bestIndex int
		var bestPred  float64
		
		Step("Building dataset", func(){
			classes := classesBuilders[argv.flags.classBuilder](ParseCommandLine(argv.flags.builderArgs))
			elements = GatherElements(classes)
			dataset = ConvertDataFrame(IndexMethod(elements), classes)
		})
		Step("Loading table", func(){
			keys = Filter(func(key string)bool{return !(key=="T"||key=="CAs"||key=="CBs"||key=="CCs")}, model.Keys)
			table = LoadMatrixTable(argv.flags.table, elements, keys)
		})
		Step("Building labels and input", func(){
			labels, input = IterateOverClasses(dataset, table, keys, concentrations, float64(600)/float64(1200))
		})
		var predictionMethod func([][]float64)[]float64
		if model.PredictColumns == nil {
			Step("Transposing", func(){
				input = Transpose(input)
			})
			predictionMethod = model.PredictRows
		} else {
			predictionMethod = model.PredictColumns
		}
		Step("Predicting", func(){
			preds = predictionMethod(input)
		})
		Step("Finding best prediction", func(){
			bestIndex = -1
			bestPred = -10
			for i, pred := range preds {
				if pred > bestPred {
					bestIndex = i
					bestPred = pred
				}
			}
		})
		fmt.Printf("elments : %v\nkeys : %v\n", func(list []string) map[int]string{result:=map[int]string{};for i,v:=range list{result[i]=v};return result}(elements), keys)
		fmt.Printf("unlabel : %v\n", UnLabel(labels[bestIndex], elements, concentrations))
		fmt.Printf("label : %v\n", labels[bestIndex])
		fmt.Printf("pred  : %v\n", bestPred)

		var bestPreds []float64
		var bestLabels []string
		var output DataFrame[string]
		Step(fmt.Sprintf("Finding %d bests", argv.flags.outputSize), func(){
			bestPreds, bestLabels = FindBest(preds, labels, argv.flags.outputSize)
		})
		var bestMaterials []string
		Step("Building material names", func(){
			bestMaterials = MultipleUnLabel(bestLabels, elements, concentrations)
		})
		Step("Building output", func(){
			output = DataFrame[string] {
				columns: map[string][]string {
					"Material": bestMaterials,
					"Prediction": Apply(float64ToString, bestPreds),
				},
			}
		})
		Step("Saving output", func(){
			output.Save(argv.flags.output)
		})
		// fmt.Printf("\n\nkeys = %v\nelements = %v\nlabels[14] = %v\ninput[14] = %v\n", keys, elements, labels[14], input[14])
	})
	
})

