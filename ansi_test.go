package cgrep

import (
	"testing"
)

func TestANSITextParse(t *testing.T) {
	colorBlue := "\x1b\x5b\x33\x34\x6d"
	colorYellow := "\x1b\x5b\x33\x33\x6d"
	colorReset := "\x1b\x5b\x30\x6d"

	text := colorBlue + "hello " + colorYellow + "world" + colorReset
	ansiText := ANSITextParse(text)
	if got, want := ansiText.Plaintext, "hello world"; got != want {
		t.Errorf("got: %v, want: %v", got, want)
	}
	ansiRanges := []ANSIEscapeCodeRange{
		{0, len("hello world") + 1, colorReset},
		{0, len("hello world"), colorBlue},
		{len("hello "), len("hello world"), colorYellow},
		{len("hello world"), len("hello world"), colorReset},
	}
	if got, want := len(ansiText.ANSIRanges), len(ansiRanges); got != want {
		t.Fatalf("got: %v, want: %v", got, want)
	}
	for i, v := range ansiText.ANSIRanges {
		if got, want := v, ansiRanges[i]; got != want {
			t.Errorf("got: %v, want: %v", got, want)
		}
	}
}
