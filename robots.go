package robots

import (
	"io"
	"net/url"
	"strings"

	"github.com/0x51-dev/robots/abnf"
	"github.com/0x51-dev/upeg/parser"
)

type File struct {
	Groups   []*Group
	Sitemaps []*url.URL
	Others   []OtherField
}

// ParseRobotsFile parses a robots.txt file from the given reader.
func ParseRobotsFile(f io.Reader) (*File, error) {
	raw, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}
	p, err := parser.New([]rune(string(raw)))
	if err != nil {
		return nil, err
	}
	n, err := p.ParseEOF(abnf.Robotstxt)
	if err != nil {
		return nil, err
	}
	return parseRobotsTXT(n)
}

func (f File) Allowed(pt ProductToken, path *url.URL) bool {
	for i := range f.Groups {
		g := f.Groups[len(f.Groups)-1-i]
		allowed := g.Allowed(pt, path)
		if allowed != nil {
			return *allowed
		}
	}

	// No rules applied.
	return true
}

func (r File) Equal(o *File) bool {
	if len(r.Groups) != len(o.Groups) {
		return false
	}
	for i := range r.Groups {
		if !r.Groups[i].Equal(o.Groups[i]) {
			return false
		}
	}
	return true
}

func (r File) String() string {
	var b strings.Builder
	for _, g := range r.Groups {
		b.WriteString(g.String())
		b.WriteString("\n")
	}
	return b.String()
}

type OtherField struct {
	Key   string
	Value string
}
