//@hts
//do not delete previous comment

package gohtslib

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"reflect"
)

func intToString(i int) string { return strconv.Itoa(i) }
func stringToInt(s string) int { result, _ := strconv.ParseInt(s, 10, 0); return int(result); }
func float64ToString(f float64) string { return fmt.Sprintf("%v", f) }
func stringToFloat64(s string) float64 { result, _ := strconv.ParseFloat(s, 64);	return float64(result); }

func AsAnySlice(obj any) ([]any, bool) {
	value := reflect.ValueOf(obj)
	if value.Kind() == reflect.Slice {
		result := make([]any, value.Len())
		for i := range result {
			result[i] = value.Index(i).Interface()
		}
		return result, true
	} else {
		return []any{}, false
	}
}

func Repr(object any) string {
	// panic("paninc")
	if cast, ok := object.(string); ok {
		return fmt.Sprintf("\"%s\"", cast)
		// TODO : handle escape chars
	} else if cast, ok := AsAnySlice(object); ok {
		return fmt.Sprintf("[%s]", StrJoin(", ", Apply(Repr, cast)))
	}
	return fmt.Sprintf("%v", object)
	// panic(fmt.Sprintf("Wrong cast : %T\n", object))
}

func StrSplit(chain string, sep string) []string {
	rchain := []rune(chain) 
	rsep := []rune(sep)
	result := []string{}
	i := 0
	j := 0
	for j=0; j<len(rchain); j++ {
		if SliceEq(rchain[j:j+len(rsep)], rsep) {
			result = append(result, string(rchain[i:j]))
			j += len(rsep)
			i = j
		}
	}
	result = append(result, string(rchain[i:j]))
	return result
}

func StrJoin(sep string, parts []string) string {
	if len(parts) == 0 {
		return ""
	} 
	if len(parts) == 1 {
		return parts[0]
	}
	result := ""
	for i, part := range parts {
		result += part
		if i < len(parts) - 1 {
			result += sep
		}
	}
	return result
}

func ReadFile(filePath string) string {
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		panic(fmt.Sprintf("Failed to read file %s: %v", filePath, err))
	}
	return string(content)
}

func WriteFile(filePath string, content string) {
	err := ioutil.WriteFile(filePath, []byte(content), 0644)
	if err != nil {
		panic(fmt.Sprintf("Failed to write to file %s: %v", filePath, err))
	}
}