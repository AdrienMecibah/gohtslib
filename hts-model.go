//@hts
//do not delete previous comment

package gohtslib

type Model struct {
	Keys []string
	PredictColumns func([][]float64)[]float64
	PredictRows func([][]float64)[]float64
}

var models map[string]Model = map[string]Model{}

func NamedModel(name string, model Model) Model {
	models[name] = model
	return model
}