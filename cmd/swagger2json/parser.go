package main

import (
	"github.com/go-openapi/spec"
)

func generateJSONSchema(inputFile, outputFolder string, pretty bool) {
	spec := getSwaggerSpec(inputFile)

	result := make(map[string]*Type)
	for k, v := range spec.Definitions {
		t := &Type{}
		t.Format = "grid-strict"
		t.Properties = make(map[string]*Type)
		t.Title = k
		t.Type = v.Type[0]

		for i, v2 := range getProperties(v) {
			t.Properties[v2.Name] = setInputAttributes(v2, getValue(i, v2, spec.Definitions))
		}

		if len(v.Enum) > 0 {
			handleEnum(t, v)
		}

		result[k] = t
	}

	doPrint(result, outputFolder, pretty)
}

func getValue(i int, o OrderedType, definitions spec.Definitions) *Type {
	t := &Type{}
	t.PropertyOrder = i + 1
	if o.Type != "" {
		t.Type = o.Type
	}
	if o.Ref != "" {
		for k, v := range definitions {
			if k == o.Ref {
				t.Type = v.Type[0]
				if v.Type[0] == "array" {
					t.Format = "tabs-top"
					t.PropertyOrder += 500
					ref := splitRef(v.Items.Schema.Ref.String())
					for k1, v1 := range definitions {
						if k1 == ref {
							t.Items = handleArray(v1, definitions)
							break
						}
					}
					break
				}
				t.Properties = make(map[string]*Type)
				for i, v2 := range getProperties(v) {
					t.Properties[v2.Name] = setInputAttributes(v2, getValue(i, v2, definitions))
				}
				if len(v.Enum) > 0 {
					handleEnum(t, v)
				}
			}
		}
	}

	if len(o.Schema.Enum) > 0 {
		handleEnum(t, o.Schema)
	}

	if o.Type == "array" {
		t.Format = "tabs-top"
		t.PropertyOrder += 500
		if o.Schema.Items.Schema.Ref.String() != "" {
			for k, v := range definitions {
				if k == splitRef(o.Schema.Items.Schema.Ref.String()) {
					t.Items = handleArray(v, definitions)
				}
			}
		} else {
			t.Items = &Type{
				Type: o.Schema.Items.Schema.Type[0],
			}
		}
	}

	if t.Type == "object" || t.Type == "array" {
		t.Options = &Options{
			GridColumns: 12,
			GridBreak:   true,
		}
	}

	return t
}
