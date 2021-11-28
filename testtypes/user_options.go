// DO NOT EDIT. This file was generated by gooptions.

package testtypes

import (
	"encoding/json"
	"time"
)

type Option func(*User)

func (u *User) with(options ...Option) *User {
	for _, option := range options {
		option(u)
	}
	return u
}

func WithEmail(email string) Option {
	return func(u *User) {
		u.email = email
	}
}

func WithFirstName(firstName string) Option {
	return func(u *User) {
		u.firstName = firstName
	}
}

func WithLastName(lastName string) Option {
	return func(u *User) {
		u.lastName = lastName
	}
}

func WithEnabled(enabled bool) Option {
	return func(u *User) {
		u.enabled = enabled
	}
}

func WithIsOrgSuperuser(isOrgSuperuser bool) Option {
	return func(u *User) {
		u.isOrgSuperuser = isOrgSuperuser
	}
}

func WithFor(for0 uintptr) Option {
	return func(u *User) {
		u.For = for0
	}
}

func WithByte(byte1 byte) Option {
	return func(u *User) {
		u.Byte = byte1
	}
}

func WithRune(rune2 rune) Option {
	return func(u *User) {
		u.Rune = rune2
	}
}

func WithCreatedBy(createdBy *string) Option {
	return func(u *User) {
		u.CreatedBy = createdBy
	}
}

func WithNumbers(numbers []int) Option {
	return func(u *User) {
		u.numbers = numbers
	}
}

func WithUuid(uuid [16]byte) Option {
	return func(u *User) {
		u.uuid = uuid
	}
}

func WithRecv(recv <-chan int) Option {
	return func(u *User) {
		u.Recv = recv
	}
}

func WithSend(send chan<- int) Option {
	return func(u *User) {
		u.Send = send
	}
}

func WithChan(chan3 chan int) Option {
	return func(u *User) {
		u.Chan = chan3
	}
}

func WithMap(map4 map[string]int) Option {
	return func(u *User) {
		u.Map = map4
	}
}

func WithF(f func(int, int, ...string) bool) Option {
	return func(u *User) {
		u.F = f
	}
}

func WithT(t time.Time) Option {
	return func(u *User) {
		u.T = t
	}
}

func WithE(e json.Encoder) Option {
	return func(u *User) {
		u.E = e
	}
}

func WithOrgs(orgs map[string]*Org) Option {
	return func(u *User) {
		u.Orgs = orgs
	}
}
