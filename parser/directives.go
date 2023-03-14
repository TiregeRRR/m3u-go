package parser

import (
	"fmt"
	"strings"
)

type ExtDirective string

func (e ExtDirective) Marshall() []byte {
	return []byte(e)
}

type Directive interface {
	Marshall() []byte
}

type EXTM3U struct{}

func (e EXTM3U) Marshall() []byte {
	return []byte("#EXTM3U\n")
}

type EXTINF struct {
	Duration  string
	TrackName string
	Path      string
}

func (e EXTINF) Marshall() []byte {
	s := fmt.Sprintf("#EXTINF:%s, %s\n%s\n", e.Duration, e.TrackName, e.Path)
	return []byte(s)
}

type PLAYLIST struct {
	Title string
}

func (p PLAYLIST) Marshall() []byte {
	s := fmt.Sprintf("#PLAYLIST:%s\n", p.Title)
	return []byte(s)
}

func parseEXTINF(lit string) EXTINF {
	titlePath := strings.SplitN(lit, ",", 2)

	duration := strings.Split(titlePath[0], ":")[1]

	titlePathSpl := strings.Split(titlePath[1], "\n")

	title := titlePathSpl[0]
	path := titlePathSpl[1]

	return EXTINF{
		Duration:  duration,
		TrackName: title,
		Path:      path,
	}
}

func parsePLAYLIST(lit string) PLAYLIST {
	title := strings.Split(lit, ":")[1]
	return PLAYLIST{Title: title}
}
