package main

import (
	"encoding/json"

	"github.com/go-openapi/spec"
)

type Schema struct {
	*Type
	Definitions Definitions `json:"definitions,omitempty"`
}

type Type struct {
	Title                string           `json:"title,omitempty"`
	Type                 string           `json:"type,omitempty"`
	Version              string           `json:"$schema,omitempty"`
	Ref                  string           `json:"$ref,omitempty"`
	MultipleOf           int              `json:"multipleOf,omitempty"`
	Maximum              int              `json:"maximum,omitempty"`
	ExclusiveMaximum     bool             `json:"exclusiveMaximum,omitempty"`
	Minimum              int              `json:"minimum,omitempty"`
	ExclusiveMinimum     bool             `json:"exclusiveMinimum,omitempty"`
	MaxLength            int              `json:"maxLength,omitempty"`
	MinLength            int              `json:"minLength,omitempty"`
	Pattern              string           `json:"pattern,omitempty"`
	AdditionalItems      *Type            `json:"additionalItems,omitempty"`
	Items                *Type            `json:"items,omitempty"`
	MaxItems             int              `json:"maxItems,omitempty"`
	MinItems             int              `json:"minItems,omitempty"`
	UniqueItems          bool             `json:"uniqueItems,omitempty"`
	MaxProperties        int              `json:"maxProperties,omitempty"`
	MinProperties        int              `json:"minProperties,omitempty"`
	Required             []string         `json:"required,omitempty"`
	Properties           map[string]*Type `json:"properties,omitempty"`
	PatternProperties    map[string]*Type `json:"patternProperties,omitempty"`
	AdditionalProperties json.RawMessage  `json:"additionalProperties,omitempty"`
	Dependencies         map[string]*Type `json:"dependencies,omitempty"`
	Enum                 []interface{}    `json:"enum,omitempty"`
	AllOf                []*Type          `json:"allOf,omitempty"`
	AnyOf                []*Type          `json:"anyOf,omitempty"`
	OneOf                []*Type          `json:"oneOf,omitempty"`
	Not                  *Type            `json:"not,omitempty"`
	Definitions          Definitions      `json:"definitions,omitempty"`
	Description          string           `json:"description,omitempty"`
	Default              interface{}      `json:"default,omitempty"`
	Format               string           `json:"format,omitempty"`
	Examples             []interface{}    `json:"examples,omitempty"`
	Media                *Type            `json:"media,omitempty"`
	BinaryEncoding       string           `json:"binaryEncoding,omitempty"`
	Options              *Options         `json:"options,omitempty"`
	PropertyOrder        int              `json:"propertyOrder,omitempty"`
}

type Definitions map[string]*Type

type Options struct {
	EnumTitles  []string `json:"enum_titles,omitempty"`
	GridColumns int      `json:"grid_columns,omitempty"`
	GridBreak   bool     `json:"grid_break,omitempty"`
}

type OrderedType struct {
	Name   string
	Schema spec.Schema
}
