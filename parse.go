package robots

import (
	"fmt"
	"net/url"
	"sort"
	"strings"

	"github.com/0x51-dev/robots/abnf"
	"github.com/0x51-dev/upeg/parser"
)

func parseIdentifier(n *parser.Node) (string, error) {
	return n.Value(), nil
}

func parsePathPattern(n *parser.Node) (*url.URL, error) {
	return url.Parse(n.Value())
}

func parseSitemap(n *parser.Node) (*url.URL, error) {
	return url.Parse(n.Children()[0].Value())
}

func parseRobotsTXT(n *parser.Node) (*File, error) {
	groups := make(map[string]*Group)

	var gs []*Group
	var sms []*url.URL
	var os []OtherField
	for _, c := range n.Children() {
		switch c.Name {
		case abnf.Group.Name:
			g, err := parseGroup(c)
			if err != nil {
				return nil, err
			}

			if len(g.ProductTokens) == 1 {
				// Only support merging groups with a single product token, for now.

				pt := g.ProductTokens[0]
				name := strings.ToLower(string(*pt))
				if group, ok := groups[name]; !ok {
					groups[name] = g
				} else {
					group.Rules = append(group.Rules, g.Rules...)
					continue // Skip adding the group to the list of groups.
				}
			}

			gs = append(gs, g)
		case abnf.Sitemap.Name:
			u, err := parseSitemap(c)
			if err != nil {
				return nil, err
			}
			sms = append(sms, u)
		case abnf.Others.Name:
			o, err := parseOther(c)
			if err != nil {
				return nil, err
			}
			os = append(os, *o)
		default:
			return nil, fmt.Errorf("unexpected root node %q", c.Name)
		}
	}
	return &File{
		Groups:   gs,
		Sitemaps: sms,
		Others:   os,
	}, nil
}

func parseGroup(n *parser.Node) (*Group, error) {
	productTokens := make(map[string]struct{})
	var pts ProductTokens
	var rs []*Rule
	for _, c := range n.Children() {
		switch c.Name {
		case abnf.Startgroupline.Name:
			pt, err := parseStartGroupLine(c)
			if err != nil {
				return nil, err
			}
			name := strings.ToLower(string(*pt))
			if _, ok := productTokens[name]; ok {
				continue
			}
			productTokens[name] = struct{}{}
			pts = append(pts, pt)
		case abnf.Rule.Name:
			r, err := parseRule(c)
			if err != nil {
				return nil, err
			}
			if r == nil {
				continue
			}
			rs = append(rs, r)
		default:
			return nil, fmt.Errorf("unexpected group node %q", c.Name)
		}
	}
	sort.Sort(pts) // Sort the product tokens.
	if *pts[0] == Wildcard {
		// The wildcard will always be the first product token if present.
		// No need to keep other product tokens if wildcard is present.
		pts = pts[:1]
	}
	return &Group{
		ProductTokens: pts,
		Rules:         rs,
	}, nil
}

func parseOther(n *parser.Node) (*OtherField, error) {
	key, err := parseIdentifier(n.Children()[0])
	if err != nil {
		return nil, err
	}
	return &OtherField{
		Key:   key,
		Value: n.Children()[1].Value(),
	}, nil
}

func parseProductToken(n *parser.Node) (*ProductToken, error) {
	if ProductToken(n.Value()) == Wildcard {
		wc := Wildcard
		return &wc, nil
	}
	id, err := parseIdentifier(n.Children()[0])
	if err != nil {
		return nil, err
	}
	pt := ProductToken(id)
	return &pt, nil
}

func parseStartGroupLine(n *parser.Node) (*ProductToken, error) {
	return parseProductToken(n.Children()[0])
}

func parseRule(n *parser.Node) (*Rule, error) {
	r, err := parseType(n.Children()[0])
	if err != nil {
		return nil, err
	}
	var pp *url.URL
	if len(n.Children()) == 2 {
		if pp, err = parsePathPattern(n.Children()[1]); err != nil {
			return nil, err
		}
	}
	return &Rule{
		Type: r,
		Path: pp,
	}, nil
}

func parseType(n *parser.Node) (RuleType, error) {
	switch n.Children()[0].Name {
	case abnf.Allow.Name:
		return Allow, nil
	case abnf.Disallow.Name:
		return Disallow, nil
	default:
		return "", fmt.Errorf("unexpected type value %q", n.Value())
	}
}
