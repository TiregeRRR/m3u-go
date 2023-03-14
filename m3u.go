package m3ugo

import (
	"bytes"
	"io"

	"github.com/TiregeRRR/m3u-go/parser"
)

type M3U struct {
	dirs []parser.Directive
}

func New(r io.Reader) (*M3U, error) {
	p := parser.NewParser(r)
	dirs, err := p.Parse()
	if err != nil {
		return nil, err
	}
	return &M3U{dirs: dirs}, nil
}

func (m *M3U) RemoveExtensions() {
	var d []parser.Directive
	for i := range m.dirs {
		switch m.dirs[i].(type) {
		case parser.ExtDirective:
			continue
		}
		d = append(d, m.dirs[i])
	}
	m.dirs = d
}

func (m *M3U) AddPrefixToPath(prefix string) {
	for i := range m.dirs {
		switch d := m.dirs[i].(type) {
		case parser.EXTINF:
			d.Path = prefix + d.Path
		}
	}
}

func (m M3U) Marshall() []byte {
	b := bytes.Buffer{}
	for i := range m.dirs {
		_, _ = b.Write(m.dirs[i].Marshall())
	}
	return b.Bytes()
}
