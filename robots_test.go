package robots_test

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/0x51-dev/robots"
)

func TestFile_Allowed(t *testing.T) {
	wc := robots.Wildcard
	block := robots.ProductToken("Block")
	f := robots.File{
		Groups: []*robots.Group{
			{
				ProductTokens: robots.ProductTokens{
					&wc,
				},
				Rules: []*robots.Rule{
					{
						Type: robots.Allow,
						Path: &url.URL{Path: "/"},
					},
				},
			},
			{
				ProductTokens: robots.ProductTokens{
					&block,
				},
				Rules: []*robots.Rule{
					{
						Type: robots.Disallow,
						Path: &url.URL{Path: "/"},
					},
				},
			},
		},
	}
	for _, test := range []struct {
		pt      string
		path    string
		allowed bool
	}{
		{"ExampleBot", "/", true},
		{"Block", "/", false},
	} {
		u, err := url.Parse(test.path)
		if err != nil {
			t.Fatal(err)
		}
		allowed := f.Allowed(robots.ProductToken(test.pt), u)
		if allowed != test.allowed {
			t.Errorf("expected %v, got %v", test.allowed, allowed)
		}
	}
}

func TestFile_Remotes(t *testing.T) {
	for _, u := range []string{
		"https://www.google.com/robots.txt",
		"https://github.com/robots.txt",
		"https://www.bing.com/robots.txt",
		"https://x.com/robots.txt",
	} {
		resp, err := http.Get(u)
		if err != nil {
			t.Fatal(err)
		}
		f, err := robots.ParseRobotsFile(resp.Body)
		if err != nil {
			t.Fatal(err)
		}
		if f == nil {
			t.Fatal("expected robots file, got nil")
		}
	}
}
