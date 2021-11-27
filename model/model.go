package model

import (
	"reflect"
)

type Model struct {
	Options *Options

	Package *Package

	StructType        *StructType
	EffectivePackages map[string]string
}

func NewModel(options *Options, p *Package, st *StructType) *Model {

	// log.Printf("Model Package: %+#v\n", *p)

	imps := st.getImports()

	// for _, imp := range imps {
	// 	log.Printf("%+#v\n", *imp)
	// }

	// log.Printf("same package field %+#v\n", st.Fields[len(st.Fields)-1].Type.getImports()[0])

	effectivePackages := CreateEffectivePackages(p, imps)
	// log.Println("ep", effectivePackages)

	return &Model{
		Options:           options,
		Package:           p,
		StructType:        st,
		EffectivePackages: effectivePackages,
	}
}

// Package path to effect import name.
// Empty effect name currently means in the same package and not to have prefix
// before dot (".").
func CreateEffectivePackages(modelPackage *Package, ps []*Package) map[string]string {
	result := map[string]string{}

	for _, p := range ps {
		if reflect.DeepEqual(modelPackage, p) {
			result[p.Path] = ""
			continue
		}
		result[p.Path] = p.Name
	}

	return result
}
