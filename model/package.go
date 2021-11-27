package model

import "strings"

type Package struct {
	Path string
	Name string
}

func NewPackage(pkgPath string) *Package {
	liSlash := strings.LastIndex(pkgPath, "/")

	return &Package{
		Path: pkgPath,
		Name: string(pkgPath[liSlash+1:]),
	}
}
