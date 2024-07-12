// Package stricks (string tricks) provides common string operations that looked like they belong here.
package stricks

import (
	"crypto/rand"
	"encoding/hex"
	"net/url"
)

func ValidURL(s string) bool {
	_, err := url.ParseRequestURI(s)
	return err == nil
}

func ParseValidURL(s string) *url.URL {
	u, err := url.ParseRequestURI(s)
	if err != nil {
		panic(err)
	}
	return u
}

func SameHost(s1, s2 string) bool {
	u1, err1 := url.ParseRequestURI(s1)
	u2, err2 := url.ParseRequestURI(s2)
	return err1 == nil && err2 == nil && u1.Host == u2.Host
}

func StringifyAnything(o any) string {
	switch s := o.(type) {
	case string:
		return s
	default:
		return ""
	}
}

func RandomWhatever() string {
	b := make([]byte, 20)
	_, err := rand.Read(b)
	if err != nil {
		return ""
	}
	return hex.EncodeToString(b[2:20])
}
