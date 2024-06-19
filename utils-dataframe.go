//@hts
//do not delete previous comment

package gohtslib

import (
	"fmt"
	"log"
	"io/ioutil"
	"sort"
)

type DataFrame[T comparable] struct {
	columns map[string][]T
	heads []string
}

func (self DataFrame[T]) Save(filename string) {
	self.Heads()
	output := StrJoin(",", self.heads)+"\n"+StrJoin("\n", Apply(
		func(row map[string]T)string{
			return StrJoin(",", Apply(
				func(head string)string{
					return fmt.Sprintf("%v", row[head])
				},
				self.heads,
			))
		},
		self.IterRows(),
	))
	WriteFile(filename, output)
}

func (self *DataFrame[T]) Length() int {
	length := -1
	for _, list := range self.columns {
		if length != -1 && length != len(list) {
			panic("Not all columns have same length")
		}
		length = len(list)
	}
	return length
}

func (self *DataFrame[T]) Width() int {
	return len(self.columns)
}

func (self *DataFrame[T]) SelectRows(cond func(map[string]T)bool) *DataFrame[T] {
	return DataFrameFromRows(Filter(cond, self.IterRows()))
}

func (self *DataFrame[T]) SelectColumns(cond func(string, []T)bool) *DataFrame[T] {
	result := DataFrame[T]{columns: make(map[string][]T, self.Width())}
	for _, head := range self.heads {
		column := self.columns[head]
		if !cond(head, column) {
			continue
		}
		newColumn := make([]T, len(column))
		for i := range column {
			newColumn[i] = column[i]
		}
		result.columns[head] = newColumn
		result.heads = append(result.heads, head)
	}
	return &result
}

func (self *DataFrame[T]) Heads() []string {
	result := make([]string, len(self.columns))
	i := 0
	for name := range self.columns {
		result[i] = name
		i += 1
	}
	_h := Copy(self.heads)
	_r := Copy(result)
	sort.Strings(_h)
	sort.Strings(_r)
	if SliceEq(_h, _r) {
		return self.heads
	}
	self.heads = result
	return result
}

func (self *DataFrame[T]) IterRows() []map[string]T {
	length := self.Length()
	result := make([]map[string]T, length)
	for i:=0; i<length; i++ {
		line := map[string]T{}
		for _, name := range self.heads {
			line[name] = self.columns[name][i]
		}
		result[i] = line
	}
	return result
}

func (self *DataFrame[T]) IterColumns() [][]T {
	result := make([][]T, len(self.columns))
	i := 0
	for _, name := range self.heads {
		result[i] = self.columns[name]
	}
	return result
}

func (self *DataFrame[T]) GetRow(i int) map[string]T {
	result := make(map[string]T, self.Width())
	for _, head := range self.heads {
		result[head] = self.columns[head][i]
	}
	return result
}

func (self *DataFrame[T]) Read(filepath string) {
    content, err := ioutil.ReadFile(filepath)
    if err != nil {
        log.Fatalf("failed to read file: %s", err)
    }
    text := string(content)
	lines := StrSplit(text, "\r\n")
	if len(lines) == 1 || len(lines) == 2 {
		lines = StrSplit(text, "\n"	)
	}
	heads := StrSplit(lines[0], ",")
	self.heads = heads
	result := map[string][]T{}
	result = make(map[string][]T, len(heads))
	for _, head := range heads {
		result[head] = make([]T, len(lines)-1)
	}
	f := func(value string) T {
		switch any(make([]T, 1)[0]).(type) {
			// case any:
			// 	i, e := strconv.ParseInt(value, 10, 0)
			// 	if e != nil {
			// 		return any(i).(T)
			// 	}
			// 	f32, e32 := strconv.ParseFloat(value, 32)
			// 	if e32 != nil {
			// 		return any(f32).(T)
			// 	}
			// 	f64, e64 := strconv.ParseFloat(value, 64)
			// 	if e64 != nil {
			// 		return any(f64).(T)
			// 	}
			// 	return make([]T, 1)[0]
			case string:
				result := any(value).(T)
				return result
			case float64:
				return any(stringToFloat64(value)).(T)
			default:
				return make([]T, 1)[0]
		}
	}
	for l, line := range lines[1:] {
		for i, value := range StrSplit(line, ",") {
			result[heads[i]][l] = f(value)
		}
	}
	self.columns = result
}
func (self *DataFrame[T]) Slice(keys ...string) *DataFrame[T] {
	content := map[string][]T{}
	for _, key := range keys {
		content[key] = Copy(self.columns[key])
	}
	result := DataFrame[T]{content, keys}
	return &result
}
func (self *DataFrame[T]) GroupBy(key string) map[T]DataFrame[T] {
	groups := []T{}
	counts := []int{}
	for _, row := range self.IterRows() {
		found := false
		for g, group := range groups {
			if group == row[key] {
				counts[g] += 1
				found = true
				break
			}
		}
		if !found {
			groups = append(groups, row[key])
			counts = append(counts, 1)
		}
	}
	result := make(map[T]DataFrame[T], len(groups))
	for g, group := range groups {
		df := DataFrame[T]{heads:self.heads}
		df.columns = make(map[string][]T, len(self.heads))
		for _, name := range self.heads {
			df.columns[name] = make([]T, counts[g])
		}
		counts[g] = 0
		for _, row := range self.IterRows() {
			if row[key] == group {
				for _, name := range self.heads {
					df.columns[name][counts[g]] = row[name]
				}
				counts[g] += 1
			}
		}
		result[group] = df
	}
	return result
}
func (self *DataFrame[T]) MultiGroupBy(keys ...string) []DataFrame[T] {
	groups := [][]T{}
	counts := []int{}
	for _, row := range self.IterRows() {
		rowGroup := make([]T, len(keys))
		for k, key := range keys {
			rowGroup[k] = row[key]
		}
		found := false
		var c int
		for g, group := range groups {
			if SliceEq(rowGroup, group) {
				found = true
				c = g
				break
			}
		}
		if !found {
			groups = append(groups, rowGroup)
			counts = append(counts, 1)
		} else {
			counts[c] += 1
		}
	}
	result := make([]DataFrame[T], len(groups))
	contents := make([]map[string][]T, len(groups))
	for c, count := range counts {
		contents[c] = map[string][]T{}
		for _, name := range self.heads {
			contents[c][name] = make([]T, count)
		}
		counts[c] = 0
	}
	for g := range groups {
		for _, row := range self.IterRows() {
			group := make([]T, len(keys))
			for k, key := range keys {
				group[k] = row[key]
			}
			if SliceEq(group, groups[g]) {
				for _, name := range self.heads {
					contents[g][name][counts[g]] = row[name]
				}
				counts[g] += 1
			}
		}
	}
	for c, content := range contents {
		result[c] = DataFrame[T]{content, self.heads}
	}
	return result
}
func (self *DataFrame[T]) AddColumn(name string, column []T) {
	if len(column) != self.Length() {
		panic("not the good size")
	}
	self.columns[name] = column
	self.heads = append(self.heads, name)
}
func (self *DataFrame[T]) String() string {
	sep := "    "
	heads := self.Heads()
	columns := make([][]string, self.Width())
	length := self.Length()
	for c := range columns {
		columns[c] = make([]string, 1+length)
		columns[c][0] = heads[c]
		for i:=0; i<length; i++ {
			columns[c][1+i] = fmt.Sprintf("%v", self.columns[heads[c]][i])
		}
	}
	for c := range columns {
		if self.Length() > 20 {
			var head []string = columns[c][:17]
			var tail []string = append([]string{"..."}, columns[c][len(columns[c])-2:]...)
			columns[c] = append(head, tail...)
			length = len(columns[c])
		}
	}
	widths := make([]int, self.Width())
	for i := range widths {
		widths[i] = -1
		for _, value := range columns[i] {
			if len(value) > widths[i] {
				widths[i] = len(value)
			}
		}
	}
	for c := range columns {
		for i, value := range columns[c] {
			prefix := ""
			for len(prefix+value) < widths[c] + 2 {
				prefix += " "
			}
			columns[c][i] = prefix + value
		}
	}
	result := ""
	for i:=0; i<length; i++ {
		if i > 0 {
			result += "\n"
		}
		line := make([]string, self.Width())
		for x := range line {
			v := columns[x][i]
			line[x] = v
		}
		result += StrJoin(sep, line)
		if i == 0 {
			bar := ""
			for j:=0; j<len(StrJoin(sep, line)); j++ {
				bar += "-"
			}
			result += "\n" + bar
		}
	}
	result += fmt.Sprintf("\n[%d rows x %d columns]\n", self.Length(), len(self.heads))
	return result
}

func ConvertDataFrame[X, Y comparable](mapper func(X)Y, dataframe *DataFrame[X]) *DataFrame[Y] {
	result := DataFrame[Y]{heads: dataframe.heads, columns: Map(
		func(name string) (string, []Y) {
			return name, Apply(
				mapper,
				dataframe.columns[name],
			)
		},
		dataframe.heads,
	)}
	return &result
}

func DataFrameFromCSV[T comparable](filepath string) *DataFrame[T] {
	result := DataFrame[T]{}
	result.Read(filepath)
	return &result
}

func DataFrameFromRows[T comparable](rows []map[string]T) *DataFrame[T] {
	result := DataFrame[T]{}
	result.heads = []string{}
	for _, row := range rows {
		for name := range row {
			found := false
			for _, head := range result.heads {
				if name == head {
					found = true
					break
				}
			}
			if !found {
				result.heads = append(result.heads, name)
			}
		}
	}
	sort.Strings(result.heads)
	result.columns = make(map[string][]T)
	for _, head := range result.heads {
		result.columns[head] = make([]T, len(rows)) 
		for r, row := range rows {
			result.columns[head][r] = row[head]
		}
	}
	return &result
}