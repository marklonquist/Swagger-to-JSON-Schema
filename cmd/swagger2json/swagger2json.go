package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"

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

	result := make(map[string]Schema)
	for k, v := range spec.Definitions {
		jsonSchema := Schema{}
		jsonSchema.Type = &Type{}
		jsonSchema.Format = "grid-strict"
		jsonSchema.Type.Properties = make(map[string]*Type)
		jsonSchema.Title = k
		jsonSchema.Type.Type = v.Type[0]

		for i, v2 := range getOrdered(v) {
			t := &Type{}
			t.PropertyOrder = i + 1

			if len(v2.Schema.Type) > 0 {
				t.Type = v2.Schema.Type[0]
				if t.Type == "array" {
					t.Format = "tabs-top"
					t.Options = &Options{
						GridColumns: 12,
						GridBreak:   true,
					}
					t.PropertyOrder += 500
					splitted := strings.Split(v2.Schema.Items.Schema.Ref.String(), "/")
					if len(splitted) == 1 {
						t.Items = &Type{
							Type: v2.Schema.Items.Schema.Type[0],
						}
					} else {
						f := splitted[2]
						t.Items = handleRef(f, spec.Definitions, i+1)
					}
				}
			} else {
				// is a $ref
				splitted := strings.Split(v2.Schema.Ref.String(), "/")
				f := splitted[2]
				t = handleRef(f, spec.Definitions, i+1)
			}

			jsonSchema.Properties[v2.Name] = t
		}

		if len(v.Enum) > 0 {
			// is enum
			handleEnum(jsonSchema.Type, v)
		}

		result[k] = jsonSchema
	}

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

func handleRef(ref string, definitions spec.Definitions, order int) *Type {
	t := &Type{}
	t.PropertyOrder = order
	t.Properties = make(map[string]*Type)

	for k, v := range definitions {
		if k == ref {
			t.Type = v.Type[0]

			if t.Type == "object" {
				t.Options = &Options{
					GridColumns: 12,
					GridBreak:   true,
				}
			}

			for i, v2 := range getOrdered(v) {
				t1 := &Type{}
				t1.PropertyOrder = i + 1
				t1.Properties = make(map[string]*Type)

				if len(v2.Schema.Type) > 0 {
					t1.Type = v2.Schema.Type[0]
					if t1.Type == "array" {
						t1.Format = "tabs-top"
						t1.Options = &Options{
							GridColumns: 12,
							GridBreak:   true,
						}
						t1.PropertyOrder += 500
						splitted := strings.Split(v2.Schema.Items.Schema.Ref.String(), "/")
						if len(splitted) == 1 {
							t1.Items = &Type{
								Type: v2.Schema.Items.Schema.Type[0],
							}
						} else {
							f := splitted[2]
							t1.Items = handleRef(f, definitions, i+1)
						}
					}
				} else {
					// is a $ref
					splitted := strings.Split(v2.Schema.Ref.String(), "/")
					f := splitted[2]
					t1 = handleRef(f, definitions, i+1)
				}

				t.Properties[v2.Name] = t1
			}

			if len(v.Enum) > 0 {
				// is enum
				handleEnum(t, v)
			}
		}
	}

	return t
}

func getOrdered(v spec.Schema) []OrderedType {
	ordered := make([]OrderedType, len(v.Properties))
	for k2, v2 := range v.Properties {
		order := int(v2.Extensions["x-position"].(float64))
		ordered[order] = OrderedType{
			Name:   k2,
			Schema: v2,
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
