package parser

import (
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/TiregeRRR/m3u-go/parser/scanner"
)

type Parser struct {
	s   *scanner.Scanner
	buf struct {
		tok scanner.Token
		lit string
		n   int
	}
}

func NewParser(r io.Reader) *Parser {
	return &Parser{s: scanner.New(r)}
}

func (p *Parser) scan() (scanner.Token, string) {
	if p.buf.n != 0 {
		p.buf.n = 0
		return p.buf.tok, p.buf.lit
	}

	tok, lit := p.s.Scan()
	p.buf.tok, p.buf.lit = tok, lit

	return tok, lit
}

func (p *Parser) unscan() { p.buf.n = 1 }

func (p *Parser) scanIgnoreWhitespace() (scanner.Token, string) {
	tok, lit := p.scan()
	if tok == scanner.WS {
		tok, lit = p.scan()
	}
	return tok, lit
}

func (p *Parser) Parse() ([]Directive, error) {
	if tok, lit := p.scanIgnoreWhitespace(); tok != scanner.DIRECTIVE && lit != "EXTM3U" {
		return nil, fmt.Errorf("found %q, expected #EXTM3U", lit)
	}
	var dirs []Directive
	dirs = append(dirs, EXTM3U{})
	for {
		tok, lit := p.scanIgnoreWhitespace()
		if tok == scanner.EOF {
			break
		}
		if tok == scanner.EXT_DIRECTIVE {
			dirs = append(dirs, ExtDirective(lit))
			continue
		}
		if tok != scanner.DIRECTIVE {
			return nil, fmt.Errorf("found %q, expected directive", lit)
		}
		dir, err := p.ParseDirective(lit)
		if err != nil {
			return nil, err
		}
		dirs = append(dirs, dir)
	}
	return dirs, nil
}

func (p *Parser) ParseDirective(lit string) (Directive, error) {
	switch {
	case strings.Contains(lit, "EXTM3U"):
		return nil, errors.New("found another declaration of EXTM3U")
	case strings.Contains(lit, "EXTINF"):
		return parseEXTINF(lit), nil
	case strings.Contains(lit, "PLAYLIST"):
		return parsePLAYLIST(lit), nil
	}
	return nil, fmt.Errorf("unsupported directive: %q", lit)
}
