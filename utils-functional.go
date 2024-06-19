//@hts
//do not delete previous comment

package gohtslib

func Apply[X, Y any](mapper func(X)Y, list []X) []Y {
	result := make([]Y, len(list))
	for i, x := range list {
		result[i] = mapper(x)
	}
	return result
}

func Filter[X any](mapper func(X)bool, list []X) []X {
	result := make([]X, 0)
	for _, x := range list {
		if mapper(x) {
			result = append(result, x)
		}
	}
	return result
}

func Map[X, V any, K comparable](mapper func(X)(K,V), list []X) map[K]V {
	result := make(map[K]V, len(list))
	for _, x := range list {
		key, value := mapper(x)
		result[key] = value
	}
	return result
}

func MapValues[K comparable, X any, Y any](mapper func(X)Y, source map[K]X) map[K]Y {
	result := make(map[K]Y, len(source))
	for key, value := range source {
		result[key] = mapper(value)
	}
	return result
}

func Range(n int) []int {
	result := make([]int, n)
	for i:=0; i<n; i++ {
		result[i] = i
	}
	return result
}