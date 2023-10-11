package sites

import (
	"github.com/pixiesys/gophetch/rules"
)

type Site interface {
	Rules() map[string]rules.Rule
	DomainKey() string
}
