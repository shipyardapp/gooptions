package main

import (
	"reflect"
)

var ReflectTypeVar = reflect.TypeOf((*struct{})(nil)).Elem()
