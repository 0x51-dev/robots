package robots

import (
	_ "embed"
	"net/url"
	"testing"

	"github.com/0x51-dev/robots/abnf"
	"github.com/0x51-dev/upeg/parser"
)

var (
	//go:embed testdata/robots.example.txt
	example string
)

func TestGroup_parse(t *testing.T) {
	for _, test := range []struct {
		input    string
		expected Group
	}{
		{
			input: "user-agent: ExampleBot\nallow: /example\n",
			expected: Group{
				ProductTokens: []*ProductToken{
					pt("ExampleBot"),
				},
				Rules: []*Rule{
					r(Allow, "/example"),
				},
			},
		},
	} {
		p, err := parser.New([]rune(test.input + "\n"))
		if err != nil {
			t.Fatal(err)
		}
		n, err := p.ParseEOF(abnf.Group)
		if err != nil {
			t.Fatal(err)
		}
		g, err := parseGroup(n)
		if err != nil {
			t.Fatal(err)
		}
		if !g.Equal(&test.expected) {
			t.Errorf("expected %v, got %v", test.expected, g)
		}
	}
}

func TestProductToken_parse(t *testing.T) {
	for _, test := range []struct {
		input    string
		expected ProductToken
	}{
		{
			input:    "user-agent: ExampleBot",
			expected: ProductToken("ExampleBot"),
		},
		{
			input:    "User-Agent: *",
			expected: Wildcard,
		},
	} {
		p, err := parser.New([]rune(test.input + "\n"))
		if err != nil {
			t.Fatal(err)
		}
		n, err := p.ParseEOF(abnf.Startgroupline)
		if err != nil {
			t.Fatal(err)
		}
		pt, err := parseStartGroupLine(n)
		if err != nil {
			t.Fatal(err)
		}
		if *pt != test.expected {
			t.Errorf("expected %q, got %q", test.expected, *pt)
		}
	}
}

func TestRobotsTXT(t *testing.T) {
	for _, test := range []struct {
		input    string
		expected File
	}{
		{
			input: "user-agent: ExampleBot\ndisallow: /foo\ndisallow: /bar\n\nuser-agent: ExampleBot\ndisallow: /baz\n",
			expected: File{
				Groups: []*Group{
					{
						ProductTokens: []*ProductToken{
							pt("ExampleBot"),
						},
						Rules: []*Rule{
							r(Disallow, "/foo"),
							r(Disallow, "/bar"),
							r(Disallow, "/baz"),
						},
					},
				},
			},
		},
		{
			input: "user-agent: ExampleBot\nuser-agent: ExampleBot\nuser-agent: *\n",
			expected: File{
				Groups: []*Group{
					{
						ProductTokens: []*ProductToken{
							pt("*"),
						},
					},
				},
			},
		},
	} {
		p, err := parser.New([]rune(test.input + "\n"))
		if err != nil {
			t.Fatal(err)
		}
		n, err := p.ParseEOF(abnf.Robotstxt)
		if err != nil {
			t.Fatal(err)
		}
		g, err := parseRobotsTXT(n)
		if err != nil {
			t.Fatal(err)
		}
		if !g.Equal(&test.expected) {
			t.Errorf("expected %v, got %v", test.expected, g)
		}
	}
}

func TestRobotsTXT_parse(t *testing.T) {
	p, err := parser.New([]rune(example))
	if err != nil {
		t.Fatal(err)
	}
	n, err := p.ParseEOF(abnf.Robotstxt)
	if err != nil {
		t.Fatal(err)
	}
	robots, err := parseRobotsTXT(n)
	if err != nil {
		t.Fatal(err)
	}
	expected := File{
		Groups: []*Group{
			{
				ProductTokens: []*ProductToken{
					pt("*"),
				},
				Rules: []*Rule{
					r(Disallow, "*.gif$"),
					r(Disallow, "/example/"),
					r(Allow, "/publications/"),
				},
			},
			{
				ProductTokens: []*ProductToken{
					pt("foobot"),
				},
				Rules: []*Rule{
					r(Disallow, "/"),
					r(Allow, "/example/page.html"),
					r(Allow, "/example/allowed.gif"),
				},
			},
			{
				ProductTokens: []*ProductToken{
					pt("barbot"),
					pt("bazbot"),
				},
				Rules: []*Rule{
					r(Disallow, "/example/page.html"),
				},
			},
			{
				ProductTokens: []*ProductToken{
					pt("examplebot"),
				},
			},
		},
	}
	if !robots.Equal(&expected) {
		t.Errorf("expected %v, got %v", expected, robots)
	}
}

func TestRule_parse(t *testing.T) {
	for _, test := range []struct {
		input    string
		expected Rule
	}{
		{
			input: "allow: /example",
			expected: Rule{
				Type: Allow,
				Path: &url.URL{Path: "/example"},
			},
		},
		{
			input: "disallow: /example",
			expected: Rule{
				Type: Disallow,
				Path: &url.URL{Path: "/example"},
			},
		},
	} {
		p, err := parser.New([]rune(test.input + "\n"))
		if err != nil {
			t.Fatal(err)
		}
		n, err := p.ParseEOF(abnf.Rule)
		if err != nil {
			t.Fatal(err)
		}
		r, err := parseRule(n)
		if err != nil {
			t.Fatal(err)
		}
		if !r.Equal(&test.expected) {
			t.Errorf("expected %v, got %v", test.expected, r)
		}
	}
}

func pt(s string) *ProductToken {
	t := ProductToken(s)
	return &t
}

func r(typ RuleType, s string) *Rule {
	u, err := url.Parse(s)
	if err != nil {
		panic(err)
	}
	return &Rule{
		Type: typ,
		Path: u,
	}
}
