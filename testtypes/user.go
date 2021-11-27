package testtypes

import (
	"encoding/json"
	"time"
)

type User struct {
	email string `gooptions:"foobar"`

	firstName string

	lastName string

	enabled bool

	isOrgSuperuser bool

	For uintptr

	Byte byte

	Rune rune

	CreatedBy *string

	numbers []int

	uuid [16]byte

	Recv <-chan int
	Send chan<- int
	Chan chan int

	Map map[string]int

	// I interface {
	// 	A() int
	// 	B(s string) bool
	// }

	F func(a int, b int, s ...string) bool

	T time.Time

	E json.Encoder

	Orgs map[string]*Org
}

type Org struct{}
