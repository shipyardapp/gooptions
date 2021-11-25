package model

import (
	"encoding/gob"
	"fmt"
	"reflect"
)

func init() {
	gob.Register(&StructType{})
	gob.Register(&StructField{})
	gob.Register(&TagOptions{})
	gob.Register(PredeclaredType(""))
}

type Type interface {
	TypeString() string
}

type StructType struct {
	Name string

	Fields []*StructField
}

func NewStructTypeFromReflectType(rt reflect.Type) (*StructType, error) {
	if rt.Kind() != reflect.Struct {
		return nil, fmt.Errorf("model: %v is not a struct type", rt)
	}

	result := &StructType{
		Name: rt.Name(),
	}

	fields, err := NewStructFieldsFromStructType(rt)
	if err != nil {
		return nil, err
	}
	result.Fields = fields

	return result, nil
}

func NewStructFieldsFromStructType(rt reflect.Type) ([]*StructField, error) {
	result := []*StructField{}

	for i := 0; i < rt.NumField(); i++ {
		sf, err := NewStructFieldFromReflectStructField(rt.Field(i))
		if err != nil {
			return nil, err
		}
		result = append(result, sf)
	}

	return result, nil
}

func NewStructFieldFromReflectStructField(sf reflect.StructField) (*StructField, error) {
	type_, err := NewStructFieldType(sf.Type)
	if err != nil {
		return nil, err
	}

	return &StructField{
		Name:       sf.Name,
		Type:       type_,
		TagOptions: &TagOptions{},
	}, nil
}

type StructField struct {
	Name string

	Type

	TagOptions *TagOptions
}

type TagOptions struct {
	Ignore bool
}

func NewStructFieldType(rt reflect.Type) (Type, error) {
	switch rt.Kind() {

	// These are all the primitive types native to Go.
	case reflect.Bool,
		reflect.String,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr,
		reflect.Float32, reflect.Float64, reflect.Complex64, reflect.Complex128:
		return PredeclaredType(rt.Kind().String()), nil
	}

	return nil, fmt.Errorf("model: unsupported struct field type %v", rt)
}

type PredeclaredType string

func (pt PredeclaredType) TypeString() string {
	return string(pt)
}
