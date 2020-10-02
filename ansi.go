package table

import (
	"fmt"
	"strconv"
)

// CodeANSI is an ANSI escape code attribute
type CodeANSI int

// ANSI colors
const (
	Black CodeANSI = iota + 30
	Red
	Green
	Yellow
	Blue
	Magenta
	Cyan
	White
)

// Bright ANSI colors
const (
	HiBlack CodeANSI = iota + 90
	HiRed
	HiGreen
	HiYellow
	HiBlue
	HiMagenta
	HiCyan
	HiWhite
)

// ANSI decorations, most not widely supported
const (
	Reset CodeANSI = iota
	Bold
	Faint
	Italic
	Underline
	BlinkSlow
	BlinkRapid
	Reverse
	Conceal
	CrossedOut
)

// Format returns a formating function applying ANSI escape codes for the given attributes.
// Please note that different terminals may support some or none of the colors and decorations.
func Format(attr ...CodeANSI) FormatFunc {
	a := buildList(attr)
	return func(s string) string {
		return fmt.Sprintf("\x1b[%sm%s\x1b[0m", a, s)
	}
}

func buildList(attr []CodeANSI) (s string) {
	for i, a := range attr {
		s += strconv.Itoa(int(a))
		if i != len(attr)-1 {
			s += ";"
		}
	}
	return
}

const (
	colorBgAdd = 10
)

// Background converts a foreground color value to it's corresponding background value
func Background(c CodeANSI) CodeANSI {
	return CodeANSI(int(c) + colorBgAdd)
}
