package text

import (
	"fmt"
	"strings"
)

const (
	Reset ColorCode = iota
	Bold
	Faint
	Italic
	Underline
	BlinkSlow
	BlinkRapid
	ReverseVideo
	Concealed
	CrossedOut
)

// Foreground text colors
const (
	FgBlack ColorCode = iota + 30
	FgRed
	FgGreen
	FgYellow
	FgBlue
	FgMagenta
	FgCyan
	FgWhite
)

type ColorCode int

// Foreground Hi-Intensity text colors
const (
	FgHiBlack ColorCode = iota + 90
	FgHiRed
	FgHiGreen
	FgHiYellow
	FgHiBlue
	FgHiMagenta
	FgHiCyan
	FgHiWhite
)

// Background text colors
const (
	BgBlack ColorCode = iota + 40
	BgRed
	BgGreen
	BgYellow
	BgBlue
	BgMagenta
	BgCyan
	BgWhite
)

// Background Hi-Intensity text colors
const (
	BgHiBlack ColorCode = iota + 100
	BgHiRed
	BgHiGreen
	BgHiYellow
	BgHiBlue
	BgHiMagenta
	BgHiCyan
	BgHiWhite
)

// ColorSting Color the given string to the given color
func ColorSting(s string, color ColorCode) string {
	return fmt.Sprintf("\033[%dm%s\033[00m", color, s)
}

func IndexOfTimes(s string, target string, times int) int {
	if times == 1 {
		return strings.Index(s, target)
	}
	place := -1
	for times > 0 {
		find := strings.Index(s, target)
		if find == -1 {
			times--
			break
		} else if find+1 >= len(s) {
			break
		}
		s = s[find+1:]
		place += find + 1
		times--
	}
	if times != 0 {
		return -1
	}
	return place
}
