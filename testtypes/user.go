package testtypes

type User struct {
	email string `gooptions:"foobar"`

	firstName string

	lastName string

	enabled bool

	isOrgSuperuser bool

	For uintptr
}
