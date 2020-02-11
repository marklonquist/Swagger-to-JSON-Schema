package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"unicode"

	"github.com/go-openapi/spec"
)

func getSwaggerSpec(inputFile string) spec.Swagger {
	file, err := os.Open(inputFile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	b, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}

	spec := spec.Swagger{}
	err = json.Unmarshal(b, &spec)
	if err != nil {
		log.Fatal(err)
	}
	return spec
}

func doPrint(result map[string]*Type, outputFolder string, pretty bool) {
	val := ""
	marshalled := make([]byte, 0)
	if pretty {
		v, err := json.MarshalIndent(result, "", "	")
		if err != nil {
			log.Fatal(err)
		}
		marshalled = v

	} else {
		v, err := json.Marshal(result)
		if err != nil {
			log.Fatal(err)
		}
		marshalled = v
	}

	val = "export const schemas = " + string(marshalled)

	outputFile, err := os.Create(outputFolder + "/jsonschema.generated.ts")
	if err != nil {
		log.Fatal(err)
	}
	defer outputFile.Close()

	outputFile.WriteString(val)
}

func handleEnum(t *Type, v spec.Schema) {
	t.Enum = make([]interface{}, len(v.Enum))
	for i, eV := range v.Enum {
		t.Enum[i], _ = eV.(float64)
	}
	if t.Options == nil {
		t.Options = &Options{}
	}
	val, _ := v.Extensions.GetStringSlice("x-enumnames")
	t.Options.EnumTitles = val
}

func setInputAttributes(v OrderedType, t *Type) *Type {
	t.Title = transformName(v.Name)

	if v.Title == "" {
		return t
	}
	if t.Options == nil {
		t.Options = &Options{
			InputAttributes: &InputAttributes{},
		}
	} else if t.Options.InputAttributes == nil {
		t.Options.InputAttributes = &InputAttributes{}
	}
	t.Options.InputAttributes.Title = v.Title
	return t
}

func transformName(name string) string {
	result := ""
	for i, r := range name {
		if i == 0 {
			r = unicode.ToUpper(r)
		} else if unicode.IsUpper(r) {
			if !unicode.IsUpper(rune(name[i-1])) {
				result += " " + string(r)
				continue
			}
		}

		result += string(r)
	}
	return result
}

func getProperties(v spec.Schema) []OrderedType {
	ordered := make([]OrderedType, len(v.Properties))
	for k2, v2 := range v.Properties {
		order := int(v2.Extensions["x-position"].(float64))
		ordered[order] = OrderedType{
			Name:   k2,
			Title:  v2.Description,
			Ref:    v2.Ref.String(),
			Schema: v2,
		}
		hidden := bool(v2.Extensions["x-hidden"].(bool))
		if hidden {
			ordered[order].Hidden = hidden
		}
		if len(v2.Type) > 0 {
			ordered[order].Type = v2.Type[0]
		}
		if ordered[order].Ref != "" {
			ordered[order].Ref = splitRef(ordered[order].Ref)
		}
	}
	return ordered
}

func splitRef(ref string) string {
	splitted := strings.Split(ref, "/")
	return splitted[2]
}

func handleArray(v spec.Schema, definitions spec.Definitions) *Type {
	t := &Type{
		Type:       v.Type[0],
		Properties: make(map[string]*Type),
	}
	for i, v2 := range getProperties(v) {
		t.Properties[v2.Name] = setInputAttributes(v2, getValue(i, v2, definitions))
	}
	return t
}
