package format

import (
	"fmt"
	"regexp"
	"strings"
)

// ParseDiscord parses an incoming Discord message to a FormattedString
func ParseDiscord(s string) FormattedString {
	return FormattedString{
		{Text: s},
	}
}

var discordEscape = regexp.MustCompile(`[\\*_]`)
var backticks = regexp.MustCompile("`+")

var trimmer = regexp.MustCompile(`^(\s*)(.*?\S)?(\s*)$`)

// RenderDiscord renders a FormattedString into a Discord message
func (fs FormattedString) RenderDiscord() string { // nolint: gocyclo
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
			} else {
				end = len(data)
			}

			if bt {
				finishedData = append(finishedData, data[begin:end]...)
			} else {
				words := strings.SplitAfter(string(data[begin:end]), " ")
				for _, word := range words {
					if strings.HasPrefix(word, "http://") || strings.HasPrefix(word, "https://") {
						finishedData = append(finishedData, []byte(word)...)
					} else {
						finishedData = append(finishedData, discordEscape.ReplaceAllFunc([]byte(word), escape)...)
					}
				}
			}
			bt = !bt
			begin = end
		}

		t := string(finishedData)

		matches := trimmer.FindAllStringSubmatch(t, -1)
		fmt.Printf("%q, %#v\n", t, matches)

		initial := matches[0][1]
		text := matches[0][2]
		final := matches[0][3]

		if text != "" {
			if (span.Format & Italic) != 0 {
				text = "*" + text + "*"
			}
			if (span.Format & Bold) != 0 {
				text = "**" + text + "**"
			}
			if (span.Format & Underline) != 0 {
				text = "__" + text + "__"
			}

			// add a U+FEFF to avoid running together format codes; "*foo*\ufeff**bar**" instead of "*foo***bar**"
			if final == "" && span.Format != None {
				text = text + "\uFEFF"
			}
		}

		output += initial + text + final
	}

	return output
}
