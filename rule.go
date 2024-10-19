package robots

import (
	"fmt"
	"net/url"
	"strings"
)

type Rule struct {
	Type RuleType
	Path *url.URL
}

// Allowed returns a pointer to a boolean indicating whether the path is allowed
// by the rule. If the path is not matched by the rule, it returns nil.
func (r Rule) Allowed(path *url.URL) *bool {
	// This is a very naive implementation of a path matching algorithm.
	matched := strings.HasPrefix(path.String(), r.Path.String())
	if !matched {
		return nil
	}
	// If the path is equal to the rule path, it's allowed.
	allowed := r.Type == Allow
	return &allowed
}

func (r Rule) Equal(o *Rule) bool {
	return r.Type == o.Type && r.Path.String() == o.Path.String()
}

func (r Rule) String() string {
	return fmt.Sprintf("%s: %s", r.Type, r.Path)
}

type RuleType string

const (
	Allow    RuleType = "allow"
	Disallow RuleType = "disallow"
)
