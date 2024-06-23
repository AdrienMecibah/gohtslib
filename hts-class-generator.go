
package gohtslib

import (
        "fmt"
)

var classesBuilders = map[string]func([]string)*DataFrame[string]{}

func GatherElements(df *DataFrame[string]) []string {
        present := Map(func(x string)(string, bool){return x, false}, CANDIDATES)
        for _, site := range SITES {
                for e, element := range df.columns[site] {
                        present[element] = true
                        found := false
                        for _, candidate := range CANDIDATES {
                                if candidate == element {
                                        found = true
                                        break
                                }
                        }
                        if !found {
                                panic(fmt.Sprintf("Unknown element %s in row #%d/%d : %s", Repr(element), e, df.Length(), Repr(df.GetRow(e))))
                        }
                }
        }
        result := []string{}
        for element, present := range present {
                if present {
                        result = append(result, element)
                }
        }
        return result
}
func NamedClassesBuilder(name string, builder func([]string)*DataFrame[string]) func([]string)*DataFrame[string] {
        _, found := classesBuilders[name]
        if found {
                panic("Class builder \""+name+"\" already exists")
        }
        classesBuilders[name] = builder
        return builder
}