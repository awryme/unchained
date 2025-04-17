package xurl

import (
	"net/url"
	"testing"
)

func TestCloneNoPanic(t *testing.T) {
	urls := []url.URL{
		{},
		{Host: "asd"},
		{Scheme: "qq", Path: "asdqwe"},
	}
	for _, u := range urls {
		Clone(&u)
	}
}
