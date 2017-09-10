package format

import (
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestParseIRC(t *testing.T) {
	Convey("When ParseIRC is used", t, func() {
		cases := []testCase{
			{"", []Span{}},
			{"foo", []Span{
				{"foo", None, Default, Default},
			}},
			{"\x0302,01foo", []Span{
				{"foo", None, Blue, Black},
			}},
			{"\x0302,01foo\x03,bar", []Span{
				{"foo", None, Blue, Black},
				{"bar", None, Blue, Default},
			}},
			{"\x0302,01foo\x03bar", []Span{
				{"foo", None, Blue, Black},
				{"bar", None, Default, Default},
			}},
			{"\x0302,01foo\x0301bar", []Span{
				{"foo", None, Blue, Black},
				{"bar", None, Black, Black},
			}},
			{"\x0302,01foo\x03,02bar", []Span{
				{"foo", None, Blue, Black},
				{"bar", None, Blue, Blue},
			}},
			{"\x02foo\x02bar\x02baz", []Span{
				{"foo", Bold, Default, Default},
				{"bar", None, Default, Default},
				{"baz", Bold, Default, Default},
			}},
			{"\x1f\x02\x1dfoo\x0fbar", []Span{
				{"foo", Bold | Italic | Underline, Default, Default},
				{"bar", None, Default, Default},
			}},
			{"\x1f\x02\x1d\x032,1foo\x0fbar", []Span{
				{"foo", Bold | Italic | Underline, Blue, Black},
				{"bar", None, Default, Default},
			}},
			{"\x1f\x02\x1d\x032,1foo\x03bar", []Span{
				{"foo", Bold | Italic | Underline, Blue, Black},
				{"bar", Bold | Italic | Underline, Default, Default},
			}},
			{"ΨΩΔ", []Span{
				{"ΨΩΔ", None, Default, Default},
			}},
		}

		for _, c := range cases {
			Convey(fmt.Sprintf("When %q is passed", c.raw), func() {
				So(ParseIRC(c.raw), ShouldResemble, c.structured)
			})
		}
	})
}

func TestRenderIRC(t *testing.T) {
	Convey("When RenderIRC is used", t, func() {
		cases := []testCase{
			{"", []Span{}},
			{"foo", []Span{
				{"foo", None, Default, Default},
			}},
			{"\x0302,01foo", []Span{
				{"foo", None, Blue, Black},
			}},
			{"\x0302,01foo\x03,bar", []Span{
				{"foo", None, Blue, Black},
				{"bar", None, Blue, Default},
			}},
			{"\x0302,01foo\x0fbar", []Span{
				{"foo", None, Blue, Black},
				{"bar", None, Default, Default},
			}},
			{"\x0302,01foo\x0301bar", []Span{
				{"foo", None, Blue, Black},
				{"bar", None, Black, Black},
			}},
			{"\x0302,01foo\x03,02bar", []Span{
				{"foo", None, Blue, Black},
				{"bar", None, Blue, Blue},
			}},
			{"\x02foo\x0fbar\x02baz", []Span{
				{"foo", Bold, Default, Default},
				{"bar", None, Default, Default},
				{"baz", Bold, Default, Default},
			}},
			{"\x02\x1d\x1ffoo\x0fbar", []Span{
				{"foo", Bold | Italic | Underline, Default, Default},
				{"bar", None, Default, Default},
			}},
			{"\x02\x1d\x1f\x0302,01foo\x0fbar", []Span{
				{"foo", Bold | Italic | Underline, Blue, Black},
				{"bar", None, Default, Default},
			}},
			{"\x02\x1d\x1f\x0302,01foo\x03bar", []Span{
				{"foo", Bold | Italic | Underline, Blue, Black},
				{"bar", Bold | Italic | Underline, Default, Default},
			}},
			{"\x02\x0302,01foo\x03\x02\x021bar", []Span{
				{"foo", Bold, Blue, Black},
				{"1bar", Bold, Default, Default},
			}},
			{"ΨΩΔ", []Span{
				{"ΨΩΔ", None, Default, Default},
			}},
		}

		for _, c := range cases {
			Convey(fmt.Sprintf("When %+v is used", c.structured), func() {
				So(c.structured.RenderIRC(), ShouldEqual, c.raw)

				// check the round-trip works too
				So(ParseIRC(c.structured.RenderIRC()), ShouldResemble, c.structured)
			})
		}
	})
}
