//@hts
//do not delete previous comment

package gohtslib

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
	"unicode"
    "unsafe"
)

type Argv[Flags any] struct {
	args []string
	present map[string]bool
	flags Flags
}

func GoifyFlag(flag string) string {
	result := ""
	runes := []rune(flag)
	for i:=0; i<len(runes); i++ {
		if runes[i] == '-' {
			i += 1
			result += strings.ToUpper(string(runes[i]))
		} else {
			result += string(runes[i])
		}
	}
	return result
}

var _ = Script("test", func(){
	inputs := []string{
		`toto titi tutu`,
		`toto "titi" tutu`,
		` toto "titi"    "tutu"    `,
		` toto 'titi'    "tutu"    `,
	}
	for _, input := range inputs {
		fmt.Printf("%s\n%v\n\n", input, ParseCommandLine(input))
	}
})

func ParseCommandLine(input string) []string {
	result := []string{}
	runes := []rune(input)
	buffer := ""
	for i:=0; i<len(runes); i++ {
		for i < len(runes) && runes[i] == ' ' {
			i += 1
		}
		if i == len(runes) {
			continue
		}
		if runes[i] == '"' {
			buffer = ""
			i += 1
			for i < len(runes) && runes[i] != '"' {
				buffer += string(runes[i])
				i += 1
			}
			result = append(result, buffer)
			i += 1
		} else if runes[i] == '\'' {
			buffer = ""
			i += 1
			for i < len(runes) && runes[i] != '\'' {
				buffer += string(runes[i])
				i += 1
			}
			result = append(result, buffer)
			i += 1
		} else {
			buffer = ""
			for i < len(runes) && runes[i] != ' ' {
				buffer += string(runes[i])
				i += 1
			}
			// i += 1
			result = append(result, buffer)			
		}
	}
	return result
}

func _ParseCommandLine(input string) []string {
	var args []string
	var current strings.Builder
	inQuotes := false
	quoteChar := rune(0)
	escaped := false
	for _, r := range input {
		switch {
			case escaped:
				current.WriteRune(r)
				escaped = false
			case r == '\\':
				escaped = true
			case r == '\'' || r == '"':
				if inQuotes && r == quoteChar {
					inQuotes = false
				} else if !inQuotes {
					inQuotes = true
					quoteChar = r
				} else {
					current.WriteRune(r)
				}
			case unicode.IsSpace(r):
				if inQuotes {
					current.WriteRune(r)
				} else if current.Len() > 0 {
					args = append(args, current.String())
					current.Reset()
				}
			default:
				current.WriteRune(r)
		}
	}
	if current.Len() > 0 {
		args = append(args, current.String())
	}
	return args
}


func ParseArgv[Flags any]() Argv[Flags] {
	argv, err := SafeParseArgv[Flags]()
	if err != nil {
		panic(err)
	}
	return argv
}
func SafeParseArgv[Flags any]() (Argv[Flags], error) {
	return SafeParseArgvFromList[Flags](os.Args[1:])
}


func ParseArgvFromLine[Flags any](line string) Argv[Flags] {
	argv, err := SafeParseArgvFromLine[Flags](line)
	if err != nil {
		panic(err)
	}
	return argv
}
func SafeParseArgvFromLine[Flags any](line string) (Argv[Flags], error) {
	return SafeParseArgvFromList[Flags](ParseCommandLine(line))
}


func ParseArgvFromList[Flags any](args []string) Argv[Flags] {
	argv, err := SafeParseArgvFromList[Flags](args)
	if err != nil {
		panic(err)
	}
	return argv
}
func SafeParseArgvFromList[Flags any](args []string) (Argv[Flags], error) {
	type Runes []rune
	positionalArguments := []string{}
	present := map[string]bool{}
	var err error
	var flags Flags
	flagsType := reflect.TypeOf(flags)
	flagsValue := reflect.ValueOf(flags)
	if flagsType.Kind() != reflect.Struct {
		return Argv[Flags]{positionalArguments, present, flags}, fmt.Errorf("provided data is not a struct")
	}
	newInstance := reflect.New(flagsType).Elem()
	structFields := make([]string, flagsValue.NumField())
	for i:=0; i<flagsValue.NumField(); i++ {
	    present[GoifyFlag(flagsType.Field(i).Name)] = false
	    structFields[i] = flagsType.Field(i).Name
	}
	for i:=0; i<len(args); i++ {
		if SliceEq(Runes(args[i])[:2], Runes("--")) {
			name := string(Runes(args[i])[2:])
			fieldName := GoifyFlag(name)
			_, found := present[fieldName]
			if !found {
				continue
				return Argv[Flags]{positionalArguments, present, flags}, fmt.Errorf("name \"%s\" is not present in struct : %v", name, structFields) 
			}
			field := newInstance.FieldByName(fieldName)
			if !field.IsValid() {
				return Argv[Flags]{positionalArguments, present, flags}, fmt.Errorf("no such field: %s in struct", fieldName)
			}
			if !field.CanSet() {
				// return Argv[Flags]{positionalArguments, present, flags}, fmt.Errorf("cannot set field: %s", fieldName)
			}
			if field.Kind() != reflect.Bool {
				return Argv[Flags]{positionalArguments, present, flags}, fmt.Errorf("field should be a boolean flag : %s", fieldName)
			}
			*(*bool)(unsafe.Pointer(field.UnsafeAddr())) = true
			present[name] = true
		} else if SliceEq(Runes(args[i])[:1], Runes("-")) {
			name := string(Runes(args[i])[1:])
			fieldName := GoifyFlag(name)
			_, found := present[fieldName]
			if !found {
				continue
				return Argv[Flags]{positionalArguments, present, flags}, fmt.Errorf("name \"%s\" is not present in struct : %v", fieldName, structFields) 
			}
			field := newInstance.FieldByName(fieldName)
			if !field.IsValid() {
				return Argv[Flags]{positionalArguments, present, flags}, fmt.Errorf("no such field: %s in struct", fieldName)
			}
			if !field.CanSet() {
				// return Argv[Flags]{positionalArguments, present, flags}, fmt.Errorf("cannot set field: %s", fieldName)
			}
			switch field.Kind() {
				case reflect.String:
					*(*string)(unsafe.Pointer(field.UnsafeAddr())) = args[i+1]
				case reflect.Int:
					intValue, err := strconv.Atoi(args[i+1])
					if err != nil {
						return Argv[Flags]{positionalArguments, present, flags}, fmt.Errorf("can not convert value from argv to type int : %s", args[i+1]) 
					}
					*(*int)(unsafe.Pointer(field.UnsafeAddr())) = intValue
				case reflect.Float32:
					floatValue, err := strconv.ParseFloat(args[i+1], 32)
					if err != nil {
						return Argv[Flags]{positionalArguments, present, flags}, fmt.Errorf("can not convert value from argv to type float32 : %s", args[i+1])
					}
					*(*float32)(unsafe.Pointer(field.UnsafeAddr())) = float32(floatValue)
				case reflect.Float64:
					floatValue, err := strconv.ParseFloat(args[i+1], 64)
					if err != nil {
						return Argv[Flags]{positionalArguments, present, flags}, fmt.Errorf("can not convert value from argv to type float64 : %s", args[i+1])
					}
					*(*float64)(unsafe.Pointer(field.UnsafeAddr())) = floatValue
				default:
					return Argv[Flags]{positionalArguments, present, flags}, fmt.Errorf("unsupported field type: %s", field.Kind())
			}
			present[GoifyFlag(name)] = true
			i += 1
		} else {
			positionalArguments = append(positionalArguments, args[i])
		}
	}
	return Argv[Flags]{positionalArguments, present, newInstance.Interface().(Flags)}, err
}