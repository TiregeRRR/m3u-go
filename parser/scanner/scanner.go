package scanner

import (
	"bufio"
	"bytes"
	"io"
	"strings"
)

type Token int

const (
	ILLEGAL Token = iota
	EOF
	WS

	DIRECTIVE
	EXT_DIRECTIVE
)

var eof = rune(0)

func isWhitespace(ch rune) bool {
	return ch == ' ' || ch == '\t' || ch == '\n'
}

func isNewline(ch rune) bool {
	return ch == '\n' || ch == '\r'
}

func isLetter(ch rune) bool {
	return (ch >= 'a' && ch <= 'z') && (ch >= 'A' && ch <= 'Z')
}

func isDigit(ch rune) bool {
	return (ch >= '0' && ch <= '9')
}

type Scanner struct {
	r *bufio.Reader
}

func New(r io.Reader) *Scanner {
	return &Scanner{r: bufio.NewReader(r)}
}

func (s *Scanner) read() rune {
	ch, _, err := s.r.ReadRune()
	if err != nil {
		return eof
	}
	return ch
}

func (s *Scanner) unread() { _ = s.r.UnreadRune() }

func (s *Scanner) Scan() (Token, string) {
	ch := s.read()

	if isWhitespace(ch) {
		s.unread()
		return s.scanWhitespace()
	} else if isLetter(ch) {
		return ILLEGAL, string(ch)
	}

	switch ch {
	case eof:
		return EOF, ""
	case '#':
		return s.scanDirective()
	}
	return ILLEGAL, string(ch)
}

func (s *Scanner) scanWhitespace() (Token, string) {
	var buf bytes.Buffer
	buf.WriteRune(s.read())
	for {
		if ch := s.read(); ch == eof {
			break
		} else if !isWhitespace(ch) {
			s.unread()
			break
		} else {
			buf.WriteRune(ch)
		}
	}
	return WS, buf.String()
}

func (s *Scanner) scanDirective() (Token, string) {
	var buf bytes.Buffer
	buf.WriteRune(s.read())

	for {
		if ch := s.read(); ch == eof {
			break
		} else if isNewline(ch) && !isDigit(ch) && ch != ':' {
			if next := s.read(); next != '#' {
				s.unread()
				buf.WriteRune('\n')
				continue
			}
			s.unread()
			break
		} else {
			buf.WriteRune(ch)
		}
	}

	str := buf.String()

	switch {
	case strings.Contains(str, "EXTM3U"):
		return DIRECTIVE, str
	case strings.Contains(str, "EXTINF"):
		return DIRECTIVE, str
	case strings.Contains(str, "PLAYLIST"):
		return DIRECTIVE, str
	case strings.Contains(str, "EXTGRP"):
		return DIRECTIVE, str
	case strings.Contains(str, "EXTALB"):
		return DIRECTIVE, str
	case strings.Contains(str, "EXTART"):
		return DIRECTIVE, str
	case strings.Contains(str, "EXTGENRE"):
		return DIRECTIVE, str
	case strings.Contains(str, "EXTM3A"):
		return DIRECTIVE, str
	case strings.Contains(str, "EXTBYT"):
		return DIRECTIVE, str
	case strings.Contains(str, "EXTBIN"):
		return DIRECTIVE, str
	case strings.Contains(str, "EXTENC"):
		return DIRECTIVE, str
	case strings.Contains(str, "EXTIMG"):
		return DIRECTIVE, str
	case strings.Contains(str, "EXT-X-"):
		return EXT_DIRECTIVE, str
	}

	return ILLEGAL, str
}
