package format

import ()

// Bitfield
type format int

// None represents no format
const None format = 0

// Format codes
const (
	Bold format = 1 << iota
	Italic
	Underline
)

// Enum
type color int

// Color codes
// These values are one more than their IRC equivalents
// This is to make the zero `color` whatever the "default" color is for the medium
const (
	Default color = iota
	White
	Black
	Blue
	Green
	BrightRed
	Red
	Magenta
	DarkYellow
	Yellow
	BrightGreen
	Cyan
	BrightCyan
	BrightBlue
	BrightMagenta
	Grey
	LightGrey

	Gray      = Grey
	LightGray = LightGrey
)

// Span represents a piece of text with a single format
type Span struct {
	Text       string
	Format     format
	Foreground color
	Background color
}

// IsZeroFormat returns whether `s` has no formatting
func (s Span) IsZeroFormat() bool {
	return s.Format == None && s.Foreground == Default && s.Background == Default
}

// IsZeroColor returns whether `s` has no color
func (s Span) IsZeroColor() bool {
	return s.Foreground == Default && s.Background == Default
}

// FormattedString represents a string made up of `Span`s
type FormattedString []Span
