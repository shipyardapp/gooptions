package model

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io"
	"log"
	"reflect"
)

func init() {
	gob.Register(&StructType{})
	gob.Register(&StructField{})
	gob.Register(&TagOptions{})
	gob.Register(PredeclaredType(""))
	gob.Register(&PointerType{})
	gob.Register(&ArraySliceType{})
	gob.Register(&ChanType{})
	gob.Register(&MapType{})
	gob.Register(&FuncType{})
	gob.Register(&NamedType{})
}

type Type interface {
	TypeString() string
}

var errorType = reflect.TypeOf((*error)(nil)).Elem()

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
	type_, err := NewType(sf.Type)
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

func NewType(rt reflect.Type) (Type, error) {
	if pkgPath := rt.PkgPath(); pkgPath != "" {
		nt := &NamedType{
			Package:       NewPackage(pkgPath),
			NameInPackage: rt.Name(),
		}
		log.Printf("%+#v %+#v\n", *nt, *nt.Package)
		return nt, nil
	}

	// Only unnamed or predeclared types after here.

	var elementType Type
	switch rt.Kind() {
	case reflect.Array, reflect.Chan, reflect.Map, reflect.Ptr, reflect.Slice:
		var err error
		elementType, err = NewType(rt.Elem())
		if err != nil {
			return nil, err
		}
	}

	switch rt.Kind() {

	// These are all the primitive types native to Go.
	case reflect.Bool,
		reflect.String,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr,
		reflect.Float32, reflect.Float64, reflect.Complex64, reflect.Complex128:
		return PredeclaredType(rt.Kind().String()), nil

	case reflect.Array:
		return &ArraySliceType{
			Len:         rt.Len(),
			ElementType: elementType,
		}, nil

	case reflect.Chan:
		return &ChanType{
			ChanDir:     rt.ChanDir(),
			ElementType: elementType,
		}, nil

	case reflect.Func:
		return NewFuncType(rt)

	case reflect.Interface:
		// Two special interface types.
		if rt.NumMethod() == 0 {
			return PredeclaredType("interface{}"), nil
		}
		if rt == errorType {
			return PredeclaredType("error"), nil
		}
		// TODO interface literal

	case reflect.Map:
		keyType, err := NewType(rt.Key())
		if err != nil {
			return nil, err
		}
		return &MapType{
			KeyType:   keyType,
			ValueType: elementType,
		}, nil

	case reflect.Ptr:
		return &PointerType{
			ElementType: elementType,
		}, nil

	case reflect.Slice:
		return &ArraySliceType{
			Len:         -1,
			ElementType: elementType,
		}, nil

	case reflect.Struct:
		// Special struct case.
		if rt.NumField() == 0 {
			return PredeclaredType("struct{}"), nil
		}
	}

	return nil, fmt.Errorf("model: unsupported type %v (%v) into model.Type", rt, rt.Kind())
}

type NamedType struct {
	Package       *Package // Could be nil.
	NameInPackage string
}

func (nt *NamedType) TypeString() string {
	return nt.Package.Name + "." + nt.NameInPackage
}

type PredeclaredType string

func (pt PredeclaredType) TypeString() string {
	return string(pt)
}

type PointerType struct {
	ElementType Type
}

func (pt *PointerType) TypeString() string {
	return "*" + pt.ElementType.TypeString()
}

type ArraySliceType struct {
	Len         int // -1 for slices, >= 0 for array types.
	ElementType Type
}

func (ast *ArraySliceType) TypeString() string {
	var l string
	if ast.Len > -1 {
		l = fmt.Sprintf("%d", ast.Len)
	}
	return fmt.Sprintf("[%v]%v", l, ast.ElementType.TypeString())
}

type ChanType struct {
	ChanDir     reflect.ChanDir
	ElementType Type
}

func (ct *ChanType) TypeString() string {
	chanString := "chan"
	switch ct.ChanDir {
	case reflect.RecvDir:
		chanString = "<-" + chanString
	case reflect.SendDir:
		chanString += "<-"
	}

	return chanString + " " + ct.ElementType.TypeString()
}

type MapType struct {
	KeyType   Type
	ValueType Type
}

func (mt *MapType) TypeString() string {
	return "map[" + mt.KeyType.TypeString() + "]" + mt.ValueType.TypeString()
}

type FuncType struct {
	In []*Parameter // The last value could have Parameter.Variadic set to true.
	// Variadic *Parameter
	Out []*Parameter
}

func NewFuncType(rt reflect.Type) (*FuncType, error) {
	nIn := rt.NumIn()
	if rt.IsVariadic() {
		nIn--
	}
	in := []*Parameter{}
	for i := 0; i < nIn; i++ {
		inType, err := NewType(rt.In(i))
		if err != nil {
			return nil, err
		}
		in = append(in, &Parameter{
			Type: inType,
		})
	}

	var variadic *Parameter
	if rt.IsVariadic() {
		variadicType, err := NewType(rt.In(nIn).Elem())
		if err != nil {
			return nil, err
		}
		variadic = &Parameter{
			Type:     variadicType,
			Variadic: true,
		}
		in = append(in, variadic)
	}

	out := []*Parameter{}
	for i := 0; i < rt.NumOut(); i++ {
		outType, err := NewType(rt.Out(i))
		if err != nil {
			return nil, err
		}
		out = append(out, &Parameter{
			Type: outType,
		})
	}

	return &FuncType{
		In: in,
		// Variadic: variadic,
		Out: out,
	}, nil
}

func (ft *FuncType) TypeString() string {
	b := &bytes.Buffer{}

	fmt.Fprint(b, "func(")
	if len(ft.In) >= 1 {
		ft.In[0].Print(b)
	}
	for i := 1; i < len(ft.In); i++ {
		fmt.Fprint(b, ", ")
		ft.In[i].Print(b)
	}
	fmt.Fprint(b, ") (")
	for _, out := range ft.Out {
		out.Print(b)
	}
	fmt.Fprint(b, ")")

	return b.String()
}

type Parameter struct {
	Name     string // Can be empty.
	Type     Type
	Variadic bool
}

func (p *Parameter) TypeString() string {
	return p.Name + p.Type.TypeString()
}

func (p *Parameter) Print(w io.Writer) {
	fmt.Fprint(w, p.Name+" ")
	if p.Variadic {
		fmt.Fprint(w, "...")
	}
	fmt.Fprint(w, p.Type.TypeString())
}
