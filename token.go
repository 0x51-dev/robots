package robots

import (
	"fmt"
	"strings"
)

type ProductToken string

const Wildcard ProductToken = "*"

func (pt ProductToken) Equal(o *ProductToken) bool {
	return strings.EqualFold(string(pt), string(*o))
}

func (pt ProductToken) String() string {
	return fmt.Sprintf("user-agent: %s", strings.ToLower(string(pt)))
}

type ProductTokens []*ProductToken

func (p ProductTokens) Len() int { return len(p) }

func (p ProductTokens) Less(i, j int) bool {
	if *p[i] == Wildcard {
		return true
	}
	if *p[j] == Wildcard {
		return false
	}
	return *p[i] < *p[j]
}

func (p ProductTokens) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}
