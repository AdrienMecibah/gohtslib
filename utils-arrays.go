//@hts
//do not delete previous comment

package gohtslib

import (
	"fmt"
)


func UpTo(n int) []int {
	result := make([]int, n)
	for i:=0; i<n; i++ {
		result[i] = i
	}
	return result
}


func Keys[K comparable, V any](mapObject map[K]V) []K {
	result := make([]K, len(mapObject))
	i := 0
	for key := range mapObject {
		result[i] = key
		i += 1
	}
	return result
}

func Values[K comparable, V any](mapObject map[K]V) []V {
	result := make([]V, len(mapObject))
	i := 0
	for _, value := range mapObject {
		result[i] = value
		i += 1
	}
	return result
}

func IsIn[T comparable](element T, list []T) bool {
	for _, e := range list {
		if element == e {
			return true
		}
	}
	return false
}

func Repeat[T any](element T, n int) []T {
	result := make([]T, n)
	for i := range result {
		result[i] = element
	}
	return result
}

func Index[T comparable](element T, list []T) int {
	for i, e := range list {
		if e == element {
			return i
		}
	}
	return -1
}

func IndexPanic[T comparable](element T, list []T) int {
	for i, e := range list {
		if e == element {
			return i
		}
	}
	panic(fmt.Sprintf("Can not find element %v in list %v", element, list))
	return -1
}

func IndexMethod[T comparable](list []T) func(T)int {
	return func(element T) int {
		return Index(element, list)
	}
}

func GetIndex[T comparable](list []T) func(int)T {
	return func(i int) T {
		return list[i] 
	}
}

func GetItem[K comparable, V any](mapping map[K]V) func(K)V {
	return func(key K) V {
		return mapping[key] 
	}
}

func TransposeMap[K1 comparable, K2 comparable, V any](mapObject map[K1]map[K2]V) map[K2]map[K1]V {
	var k2Keys []K2
	for k1 := range mapObject {
		keys := Keys(mapObject[k1])
		if k2Keys == nil {
			k2Keys = keys
		} else {
			msg := fmt.Sprintf("Not all key sets within are the same : %v and %v", k2Keys, keys)
			if len(k2Keys) != len(keys) {
				panic(msg)
			}
			for _, k2 := range k2Keys {
				found := false
				for _, k := range keys {
					if k == k2 {
						found = true
					}
				}
				if !found {
					panic(msg)
				}
			}
		}
	}
	result := make(map[K2]map[K1]V, len(k2Keys))
	for _, k2 := range k2Keys {
		line := make(map[K1]V, len(mapObject))
		for k1 := range mapObject {
			line[k1] = mapObject[k1][k2]
		}
		result[k2] = line
	}
	return result
}

func Transpose[T any](matrix [][]T) [][]T {
	// fmt.Printf("\x1b[92m%v\x1b[39m\n", matrix)
	length := -1
	for _, line := range matrix {
		if !(!(length!=-1) || length==len(line)) {
			panic(fmt.Sprintf(
				"not all of the %d columns are the same length : %v", 
				len(matrix),
				Apply(
					func(x []T) int {
						return len(x)
					}, 
					matrix,
				),
			))
		}
		if length == -1 {
			length = len(line)
		}
	}
	result := make([][]T, len(matrix[0]))
	for i := range result {
		result[i] = make([]T, len(matrix))
	}
	for i := range matrix {
		for j := range matrix[i] {
			result[j][i] = matrix[i][j]
		}
	}
	return result
}

func SliceEq[T comparable](s1, s2 []T) bool {
	if len(s1) != len(s2) {
		return false
	}
	for i := range s1 {
		if s1[i] != s2[i] {
			return false
		}
	}
	return true
}

func SetEq[T comparable](s1, s2 []T) bool {
	if len(s1) != len(s2) {
		return false
	}
	for _, x1 := range s1 {
		found := false
		for _, x2 := range s2 {
			if x1 == x2 {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}

func Copy[T any](array []T) []T {
	result := make([]T, len(array))
	for i, value := range array {
		result[i] = value
	}
	return result
}

func WithoutDuplicates[T comparable](list []T) []T {
	result := []T{}
	for i, value := range list {
		found := false
		for _, v := range list[:i] {
			if value == v {
				found = true
				break
			}
		}
		if !found {
			result = append(result, value)
		}
	}
	return result
}

func Product[T any](permutations ...[]T) [][]T {
	size := 1
	for _, permutation := range permutations {
		size *= len(permutation)
	}
	index := 0
	result := make([][]T, size)
	var f func([]T, [][]T)
	f = func(base []T, remaining[][]T) {
		if len(remaining) == 0 {
			result[index] = base
			index += 1
		} else {
			for _, value := range remaining[0] {
				f(append(base, value), remaining[1:])
			}
		}
	}
	f([]T{}, permutations)
	return result
}