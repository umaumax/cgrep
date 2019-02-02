package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/mgutz/ansi"
)

var (
	verbose bool
)

func init() {
	flag.BoolVar(&verbose, "verbose", false, "verbose flag")
}

func main() {
	flag.Parse()
	args := flag.Args()
	if flag.NArg() == 0 {
		log.Fatalf("[REG pattern with ()] expected e.g. '([0-9]+)\\.([0-9]+)f'\n")
		return
	}
	pattern := args[0]
	colorListText := ""
	if flag.NArg() > 1 {
		colorListText = args[1]
	}

	colorReg := regexp.MustCompilePOSIX(pattern)
	ansiContSeqReg := regexp.MustCompilePOSIX(`(\x9B|\x1B\[)[0-?]*[ -\/]*[@-~]`)

	defaultColorTable := []string{
		ansi.Green, ansi.Yellow, ansi.Cyan, ansi.Magenta, ansi.Blue, ansi.Red,
		ansi.LightGreen, ansi.LightYellow, ansi.LightCyan, ansi.LightMagenta, ansi.LightBlue, ansi.LightRed,
	}
	// FYI: see color format
	// [mgutz/ansi: Small, fast library to create ANSI colored strings and codes\. \[go, golang\]]( https://github.com/mgutz/ansi )
	colorTable := defaultColorTable
	if colorListText != "" {
		colorList := strings.Split(colorListText, ",")
		colorTable = make([]string, len(colorList))
		notFoundColorCode := ansi.ColorCode("NOT_FOUND_DUMMY_COLOR_CODE")
		for i, v := range colorList {
			colorCode := ansi.ColorCode(v)
			if v != "" && colorCode == notFoundColorCode {
				log.Fatalf("invalid color text:[%s]\n", v)
			}
			colorTable[i] = colorCode
		}
	}

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		text := scanner.Text()

		type ANSI_EscapeCodeRange struct {
			start int
			end   int
			code  string
		}
		ansiCodes := ansiContSeqReg.FindAllStringIndex(text, -1)
		plaintext := ansiContSeqReg.ReplaceAllString(text, "")
		plaintextRunes := []rune(plaintext)
		lenRunePlaintext := len(plaintextRunes)
		ansiRanges := make([]ANSI_EscapeCodeRange, len(ansiCodes)+1)
		// NOTE: reset per line
		ansiRanges[0] = ANSI_EscapeCodeRange{0, lenRunePlaintext + 1, ansi.Reset}

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
			ansiRanges[i+1] = ANSI_EscapeCodeRange{lenPreIndex, lenRunePlaintext, string(code)}
			preIndex = end
		}

		if verbose {
			fmt.Println("[START]")
			fmt.Println("[rawtext]", text)
			fmt.Println("[plaintext]", plaintext)
			fmt.Println("[INFO] ansi len", len(text), "plaintext string len", len(plaintext), "rune len", len([]rune(plaintext)))
		}

		m := colorReg.FindAllStringSubmatchIndex(plaintext, -1)
		if len(m) > 0 {
			for _, v := range m {
				// NOTE: [0],[1]: entire index set
				for i := 1; i*2 < len(v); i++ {
					start := v[i*2+0]
					end := v[i*2+1]
					// NOTE: not hit () e.g. '(a)|(b)'
					if start < 0 || end < 0 {
						continue
					}
					colorIndex := i - 1
					code := colorTable[colorIndex%len(colorTable)]
					// NOTE: skip coloring
					if code == "" {
						continue
					}
					ansi := ANSI_EscapeCodeRange{len([]rune(plaintext[:start])), len([]rune(plaintext[:end])), string(code)}
					ansiRanges = append(ansiRanges, ansi)
				}
			}
		}

		// NOTE: only for debug
		// for i, v := range ansiRanges {
		// fmt.Printf("%2d: %2d-%2d %q\n", i, v.start, v.end, v.code)
		// }
		codes := make([]string, lenRunePlaintext+1)
		for _, v := range ansiRanges {
			for i := v.start; i < v.end; i++ {
				codes[i] = v.code
			}
		}
		if verbose {
			fmt.Println("[OUTPUT]")
		}
		buf := new(bytes.Buffer)
		for i := 0; i < lenRunePlaintext+1; i++ {
			if codes[i] != "" {
				fmt.Fprintf(buf, "%s", codes[i])
			}
			if i == lenRunePlaintext {
				break
			}
			fmt.Fprintf(buf, "%c", plaintextRunes[i])
		}
		fmt.Fprintf(buf, "\n")
		fmt.Printf("%s", buf.String())
		if verbose {
			fmt.Println("[END]")
		}
	}
	if err := scanner.Err(); err != nil {
		log.Printf("stdin read err:%s\n", err)
	}
}
