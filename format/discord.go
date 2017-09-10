package format

import (
	"fmt"
	"regexp"
)

// ParseDiscord parses an incoming Discord message to a FormattedString
func ParseDiscord(s string) FormattedString {
	return FormattedString{
		{Text: s},
	}
}

var discordEscape = regexp.MustCompile(`\*|_`)
var backticks = regexp.MustCompile("`+")

// RenderDiscord renders a FormattedString into a Discord message
func (fs FormattedString) RenderDiscord() string {
	output := ""
	escape := func(s []byte) []byte { return append([]byte("\\"), s...) }
	for _, span := range fs {

		finishedData := []byte{}
		data := []byte(span.Text)

		begin := 0
		var end int
		bt := false
		for begin < len(data) {
			firstBackticks := backticks.FindIndex(data[begin:])
			if firstBackticks != nil {
				end = begin + firstBackticks[1]
				fmt.Printf("%q %q %q %+v %+v\n", string(finishedData), string(data), string(data[begin:]), begin, end)
				bt = !bt
			} else {
				end = len(data)
			}

			if !bt {
				finishedData = append(finishedData, data[begin:end]...)
			} else {
				finishedData = append(finishedData, discordEscape.ReplaceAllFunc(data[begin:end], escape)...)
			}
			begin = end
		}

		t := string(finishedData)

		if (span.Format & Italic) != 0 {
			t = "*" + t + "*"
		}
		if (span.Format & Bold) != 0 {
			t = "**" + t + "**"
		}
		if (span.Format & Underline) != 0 {
			t = "__" + t + "__"
		}
		output += t + "\ufeff" // add a U+FEFF to avoid running together format codes; "*foo*\ufeff**bar**" instead of "*foo***bar**"
	}

	return output
}
