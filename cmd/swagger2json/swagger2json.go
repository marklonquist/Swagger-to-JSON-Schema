package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"unicode"

	"github.com/go-openapi/spec"
	"github.com/urfave/cli"
)

var app = cli.NewApp()

func info() {
	app.Name = "JSON Schema Generator CLI"
	app.Usage = ""
	app.UsageText = "A CLI used to generate JSON Schema based on swagger:model comments"
	app.Author = "Mark Bøg Lønquist"
	app.Version = "1.0.0"
}

func commands() {
	app.Commands = []cli.Command{
		{
			Name:  "generate",
			Usage: "",
			Action: func(c *cli.Context) {
				inputFile := c.Args().Get(0)
				outputFolder := c.Args().Get(1)
				pretty := c.Args().Get(2)
				if inputFile == "" {
					log.Println("Missing input file - should point to a swagger file")
					return
				}
				if outputFolder == "" {
					log.Println("Missing output folder (should point to a folder, overwrites existing files)")
					return
				}

				generateJsonSchema(inputFile, outputFolder, pretty == "true")
			},
		},
	}
}

func main() {
	info()
	commands()

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func generateJsonSchema(inputFile, outputFolder string, pretty bool) {
	spec := getSwaggerSpec(inputFile)

	result := make(map[string]Schema)
	for k, v := range spec.Definitions {
		jsonSchema := Schema{}
		jsonSchema.Type = &Type{}
		jsonSchema.Format = "grid-strict"
		jsonSchema.Type.Properties = make(map[string]*Type)
		jsonSchema.Title = k
		jsonSchema.Type.Type = v.Type[0]

		for i, v2 := range getOrdered(v) {
			jsonSchema.Properties[v2.Name] = setInputAttributes(v2, getValue(i, v2, spec.Definitions, ""))
		}

		if len(v.Enum) > 0 {
			handleEnum(jsonSchema.Type, v)
		}

		result[k] = jsonSchema
	}

	doPrint(pretty, result, outputFolder)
}

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

func doPrint(pretty bool, result map[string]Schema, outputFolder string) {
	val := ""
	if pretty {
		marshalled, err := json.MarshalIndent(result, "", "	")
		if err != nil {
			log.Fatal("5", err)
		}
		val = "export const schemas = " + string(marshalled)
	} else {
		marshalled, err := json.Marshal(result)
		if err != nil {
			log.Fatal("5", err)
		}
		val = "export const schemas = " + string(marshalled)
	}

	outputFile, err := os.Create(outputFolder + "/jsonschema.generated.ts")
	if err != nil {
		log.Fatal(err)
	}
	defer outputFile.Close()

	outputFile.WriteString(val)
}

func getValue(i int, v *OrderedType, definitions spec.Definitions, ref string) *Type {
	t := &Type{}
	t.PropertyOrder = i + 1

	if ref != "" {
		t.Properties = make(map[string]*Type)
		for k, v2 := range definitions {
			if k == ref {
				t.Type = v2.Type[0]

				if t.Type == "array" {
					// check for refs
					splitted := strings.Split(v2.Items.Schema.Ref.String(), "/")
					f := splitted[2]
					t.Items = getValue(i, nil, definitions, f)
					break
				}

				if t.Type == "object" {
					t.Options = &Options{
						GridColumns: 12,
						GridBreak:   true,
					}
				}

				for i, v3 := range getOrdered(v2) {
					t.Properties[v3.Name] = setInputAttributes(v3, getValue(i, v3, definitions, ""))
				}

				if len(v2.Enum) > 0 {
					handleEnum(t, v2)
				}
				break
			}
		}
	} else {
		if len(v.Schema.Type) > 0 {
			t.Type = v.Schema.Type[0]
			if t.Type == "array" {
				t.Format = "tabs-top"
				t.Options = &Options{
					GridColumns: 12,
					GridBreak:   true,
				}
				t.PropertyOrder += 500
				splitted := strings.Split(v.Schema.Items.Schema.Ref.String(), "/")
				if len(splitted) == 1 {
					t.Items = &Type{
						Type: v.Schema.Items.Schema.Type[0],
					}
				} else {
					f := splitted[2]
					t.Items = getValue(i, nil, definitions, f)
				}
			}
		} else {
			// is a $ref
			splitted := strings.Split(v.Schema.Ref.String(), "/")
			f := splitted[2]
			t = getValue(i, nil, definitions, f)
		}
	}

	return t
}

func getOrdered(v spec.Schema) []*OrderedType {
	ordered := make([]*OrderedType, len(v.Properties))
	for k2, v2 := range v.Properties {
		order := int(v2.Extensions["x-position"].(float64))
		ordered[order] = &OrderedType{
			Name:   k2,
			Title:  v2.Description,
			Schema: &v2,
		}
	}
	return ordered
}

func handleEnum(t *Type, v spec.Schema) {
	t.Enum = make([]interface{}, len(v.Enum))
	for i, eV := range v.Enum {
		t.Enum[i], _ = strconv.Atoi(eV.(string))
	}
	t.Options = &Options{}
	val, _ := v.Extensions.GetStringSlice("x-enumnames")
	t.Options.EnumTitles = val
}

func setInputAttributes(v *OrderedType, t *Type) *Type {
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
