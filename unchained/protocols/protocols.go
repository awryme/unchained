package protocols

import (
	"fmt"
	"strings"
)

const Trojan = "trojan"
const Vless = "vless"

var List = []string{Trojan, Vless}

func ErrInvalid(proto string) error {
	return fmt.Errorf(
		"invalid protocol %s (options = %s)",
		proto,
		strings.Join(List, "/"),
	)
}
