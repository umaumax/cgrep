package cgrep

import (
	"bytes"
	"fmt"
	"regexp"

	"github.com/mgutz/ansi"
)

var ansiContSeqReg = regexp.MustCompilePOSIX(`(\x9B|\x1B\[)[0-?]*[ -\/]*[@-~]`)

type ANSIEscapeCodeRange struct {
	Start int
	End   int
	Code  string
}

type ANSIText struct {
	text           string
	ANSICodes      [][]int
	Plaintext      string
	PlaintextRunes []rune
	ANSIRanges     []ANSIEscapeCodeRange
}

func (a *ANSIText) PlainTextRunesLen() int {
	return len(a.PlaintextRunes)
}

func (a *ANSIText) Debug() {
	fmt.Println("[text]", a.text)
	fmt.Println("[plaintext]", a.Plaintext)
	fmt.Println("[INFO] input text len", len(a.text), "plaintext string len", len(a.Plaintext), "plaintext rune len", a.PlainTextRunesLen())
}

func (a *ANSIText) DebugANSIRanges() {
	for i, v := range a.ANSIRanges {
		fmt.Printf("%2d: %2d-%2d %q\n", i, v.Start, v.End, v.Code)
	}
}

func (a *ANSIText) String() string {
	codes := make([]string, a.PlainTextRunesLen()+1)
	for _, v := range a.ANSIRanges {
		for i := v.Start; i < v.End; i++ {
			codes[i] = v.Code
		}
	}
	buf := new(bytes.Buffer)
	for i := 0; i < len(codes); i++ {
		if codes[i] != "" {
			fmt.Fprintf(buf, "%s", codes[i])
		}
		if i < len(codes)-1 {
			fmt.Fprintf(buf, "%c", a.PlaintextRunes[i])
		}
	}
	return buf.String()
}

func ANSITextParse(text string) *ANSIText {
	ansiCodes := ansiContSeqReg.FindAllStringIndex(text, -1)
	plaintext := ansiContSeqReg.ReplaceAllString(text, "")
	plaintextRunes := []rune(plaintext)
	lenRunePlaintext := len(plaintextRunes)
	ansiRanges := make([]ANSIEscapeCodeRange, len(ansiCodes)+1)
	// NOTE: reset per line
	ansiRanges[0] = ANSIEscapeCodeRange{0, lenRunePlaintext + 1, ansi.Reset}
	// NOTE:
	lenPreIndex := 0
	preIndex := 0
	for i, v := range ansiCodes {
		// NOTE: byte index
		start := v[0]
		end := v[1]
		// NOTE: ansi以外の文字のrune indexを蓄積
		lenPreIndex += len([]rune(text[preIndex:start]))
		// NOTE: ansi color code
		code := text[start:end]
		ansiRanges[i+1] = ANSIEscapeCodeRange{lenPreIndex, lenRunePlaintext, string(code)}
		preIndex = end
	}
	return &ANSIText{
		text:           text,
		ANSICodes:      ansiCodes,
		Plaintext:      plaintext,
		PlaintextRunes: plaintextRunes,
		ANSIRanges:     ansiRanges,
	}
}
