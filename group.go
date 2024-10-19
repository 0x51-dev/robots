package robots

import (
	"net/url"
	"strings"
)

type Group struct {
	ProductTokens ProductTokens
	Rules         []*Rule
}

// Allowed returns (true, _) if the group is not for the given user-agent.
// Otherwise, it returns (_, true) if the path is allowed.
func (g Group) Allowed(pt ProductToken, path *url.URL) *bool {
	var check bool
	for _, gpt := range g.ProductTokens {
		if *gpt == Wildcard || gpt.Equal(&pt) {
			check = true
			break
		}
	}
	if !check {
		// Not checking this group, it's not for this user-agent.
		return nil
	}

	for _, r := range g.Rules {
		allowed := r.Allowed(path)
		if allowed == nil {
			// Rule path does not match the given path.
			continue
		}
		return allowed
	}

	// It matched the group, but no rules applied.
	defaultAllow := true
	return &defaultAllow
}

func (g Group) Equal(o *Group) bool {
	if len(g.ProductTokens) != len(o.ProductTokens) {
		return false
	}
	if len(g.Rules) != len(o.Rules) {
		return false
	}
	for i := range g.ProductTokens {
		if !g.ProductTokens[i].Equal(o.ProductTokens[i]) {
			return false
		}
	}
	for i := range g.Rules {
		if !g.Rules[i].Equal(o.Rules[i]) {
			return false
		}
	}
	return true
}

func (g Group) String() string {
	var b strings.Builder
	for _, pt := range g.ProductTokens {
		b.WriteString(pt.String())
		b.WriteString("\n")
	}
	for _, r := range g.Rules {
		b.WriteString(r.String())
		b.WriteString("\n")
	}
	return b.String()
}
