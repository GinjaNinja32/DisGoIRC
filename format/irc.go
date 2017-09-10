package format

import (
	"fmt"
	"regexp"
	"strconv"
)

const (
	ircColor = "\x03"

	ircReset     = "\x0f"
	ircBold      = "\x02"
	ircUnderline = "\x1f"
	ircItalic    = "\x1d"
)

var ircFormatChars = map[byte]format{
	'\x02': Bold,
	'\x1f': Underline,
	'\x1d': Italic,
}

// ParseIRC parses an incoming IRC message into a FormattedString
func ParseIRC(s string) FormattedString {
	spans := []Span{}

	currentSpan := Span{}

	data := []byte(s)

	for i := 0; i < len(data); i++ {
		c := data[i]

		if _, ok := ircFormatChars[c]; !ok && c != '\x03' && c != '\x0f' {
			currentSpan.Text += string(c)
			continue
		}

		if currentSpan.Text != "" {
			spans = append(spans, currentSpan)
		}

		currentSpan.Text = ""

		if formatCode, ok := ircFormatChars[c]; ok {
			currentSpan.Format ^= formatCode
		} else if c == '\x0f' {
			currentSpan.Format = None
			currentSpan.Foreground = Default
			currentSpan.Background = Default
		} else if c == '\x03' {
			i += parseIRCColorCode(data[i+1:], &currentSpan)
		} else {
			panic("invalid code reached")
		}
	}

	if currentSpan.Text != "" {
		spans = append(spans, currentSpan)
	}

	return FormattedString(spans)
}

var ircColorCode = regexp.MustCompile("([0-9]{0,2})(,([0-9]{0,2}))?")

func parseIRCColorCode(data []byte, currentSpan *Span) (length int) {
	colorCode := ircColorCode.FindSubmatch(data)

	fgSpecifier := -1
	if len(colorCode[1]) != 0 {
		fgSpecifier, _ = strconv.Atoi(string(colorCode[1]))
	}

	bgSpecifier := -1
	if len(colorCode[3]) != 0 {
		bgSpecifier, _ = strconv.Atoi(string(colorCode[3]))
	}

	if len(colorCode[2]) != 0 { // If comma was found
		if fgSpecifier != -1 {
			currentSpan.Foreground = color(fgSpecifier + 1)
		}

		currentSpan.Background = color(bgSpecifier + 1)
	} else {
		if fgSpecifier == -1 {
			currentSpan.Background = Default
		}
		currentSpan.Foreground = color(fgSpecifier + 1)
	}

	return len(colorCode[0])
}

// RenderIRC renders a FormattedString into an IRC message
func (fs FormattedString) RenderIRC() string { // nolint: gocyclo
	output := ""

	var lastSpan Span
	for _, span := range fs {
		if span.IsZeroFormat() && !lastSpan.IsZeroFormat() {
			output += ircReset + span.Text
			lastSpan = span
			continue
		}

		formatChanges := span.Format ^ lastSpan.Format
		if (formatChanges & Bold) != 0 {
			output += ircBold
		}
		if (formatChanges & Italic) != 0 {
			output += ircItalic
		}
		if (formatChanges & Underline) != 0 {
			output += ircUnderline
		}

		if span.IsZeroColor() && !lastSpan.IsZeroColor() {
			output += ircColor

			if !isByteSafeAfterIncompleteColor(span.Text[0]) {
				output += ircBold + ircBold
			}
		} else if span.Foreground != lastSpan.Foreground && span.Background != lastSpan.Background {
			fgSpecifier := colorToSpecifier(span.Foreground)
			bgSpecifier := colorToSpecifier(span.Background)

			output += ircColor + fgSpecifier + "," + bgSpecifier
		} else if span.Foreground != lastSpan.Foreground {
			fgSpecifier := colorToSpecifier(span.Foreground)

			output += ircColor + fgSpecifier
		} else if span.Background != lastSpan.Background {
			bgSpecifier := colorToSpecifier(span.Background)

			output += ircColor + "," + bgSpecifier
		}

		output += span.Text
		lastSpan = span
	}

	return output
}

func colorToSpecifier(c color) string {
	if c == Default {
		return ""
	}
	return fmt.Sprintf("%02d", c-1)
}

func isByteSafeAfterIncompleteColor(c byte) bool {
	return !(('0' <= c && c <= '9') || c == ',')
}
