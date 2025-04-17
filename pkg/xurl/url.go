package xurl

import (
	"fmt"
	"net/url"
)

func Clone(u *url.URL) *url.URL {
	clone, err := url.Parse(u.String())
	if err != nil {
		// cloning a url should never result in error
		panic(fmt.Sprint("clone url:", err))
	}
	return clone
}

func With(u *url.URL, updates ...func(*url.URL)) *url.URL {
	clone := Clone(u)
	for _, upd := range updates {
		upd(clone)
	}
	return clone
}
