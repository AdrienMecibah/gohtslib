//@hts
//do not delete previous comment

package gohtslib

import (
	"fmt"
	"time"
)

import (
	"math"
	"sort"
	"math/rand"
)

func Step(name string, code func()()) {
	AsStep(name, func()int{code(); return 0})
}

func AsStep[T any](name string, code func()T) T {
	start := time.Now()
	fmt.Printf("\x1b[38;5;226m%s...\x1b[39m\n", name)
	result := code()
	duration := time.Now().Sub(start)
	fmt.Printf("\x1b[38;5;226m%s : %s\x1b[39m\n", name, duration)
	return result
}

func Rmse(target []float64, pred []float64) float64 {
	if len(target) != len(pred) {
		panic("target and prediction have different size")
	}
	squaredSum := float64(0)
	length := float64(len(target))
	for i := range target {
		fmt.Printf("[%d] %v %v \n", i, target[i], pred[i])
		if target[i] == 0 && fmt.Sprintf("%v", pred[i]) == "NaN" {
			length -= 1
			continue
		}
		squaredSum += math.Pow(target[i]-pred[i], 2)
	}
	println("squaredSum = ", squaredSum, "length =", length, "len(target) =", len(target))
	return math.Pow(squaredSum/length, 0.5)
}

func _FindBest(values []float64, labels []string, n int) ([]float64, []string) {
	type Pair struct {
		Value float64
		Label string
	}
	if len(values) != len(labels) {
		panic("values and labels must have the same length")
	}

	// Create a slice of pairs
	pairs := make([]Pair, len(values))
	for i := range values {
		pairs[i] = Pair{Value: values[i], Label: labels[i]}
	}

	// Sort pairs by Value in descending order
	sort.Slice(pairs, func(i, j int) bool {
		return pairs[i].Value > pairs[j].Value
	})

	// Keep the top n pairs
	if n > len(pairs) {
		n = len(pairs)
	}
	topPairs := pairs[:n]

	// Separate the pairs back into two lists
	topValues := make([]float64, n)
	topLabels := make([]string, n)
	for i := range topPairs {
		topValues[i] = topPairs[i].Value
		topLabels[i] = topPairs[i].Label
	}

	return topValues, topLabels
}

var _ = Script("test", func(){
	x1 := Apply(func(i int)float64{return float64(i)}, UpTo(20))
	y1 := Apply(float64ToString, x1)
	rand.Shuffle(len(x1), func(i, j int) {
		x1[i], x1[j] = x1[j], x1[i]
	})
	x2, y2 := FindBest(x1, y1, 10)
	fmt.Printf("x1 = %v\ny1 = %v\nx2 = %v\ny2 = %v\n", x1, y1, x2, y2)
})

func FindBest(values []float64, labels []string, n int) ([]float64, []string) {
	type Pair struct {
		value float64
		label string
	}
	if len(values) != len(labels) {
		panic("FindBest : values and labels must have the same length")
	}
	if len(values) < n {
		return values, labels
	}
	pairs := make([]Pair, len(values))
	for i := range values {
		pairs[i] = Pair{value: values[i], label: labels[i]}
	}
	topPairs := make([]Pair, n)
	for i := range topPairs {
		topPairs[i] = pairs[i]
	}
	sort.Slice(topPairs, func(i, j int) bool {
		return topPairs[i].value > topPairs[j].value
	})
	offset := func(topPairs []Pair, value Pair, index int) {
		for i:=n-1; i!=index; i-- {
			topPairs[i] = topPairs[i-1]
		}
		topPairs[index] = value
	}
	for i:=n; i<len(values); i++ {
		if values[i] >= topPairs[n-1].value {
			j := n-1
			for j >= 0 && values[i] >= topPairs[j].value {
				j += -1
			}
			offset(topPairs, Pair{values[i], labels[i]}, j+1)
		}
	}
	topValues := make([]float64, n)
	topLabels := make([]string, n)
	for i := range topPairs {
		topValues[i] = topPairs[i].value
		topLabels[i] = topPairs[i].label
	}
	return topValues, topLabels
}