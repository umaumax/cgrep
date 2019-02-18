package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/mgutz/ansi"
	"github.com/umaumax/cgrep"
)

var (
	verbose        bool
	checkRegexFlag bool
	fixedFlag      bool
)

func init() {
	flag.BoolVar(&verbose, "verbose", false, "verbose flag")
	flag.BoolVar(&checkRegexFlag, "n", false, "check regex only")
	flag.BoolVar(&fixedFlag, "F", false, "use fixed string or not")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, ""+
			os.Args[0]+` [REG pattern with ()]
  e.g.
    cgrep '([0-9]+)\.([0-9]+)f()' 'green,default'

`)
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, `
* color format
  * See: https://github.com/mgutz/ansi
`)
	}
}

func main() {
	flag.Parse()
	args := flag.Args()
	if flag.NArg() == 0 {
		flag.Usage()
		os.Exit(1)
	}
	pattern := args[0]
	colorListText := ""
	if flag.NArg() > 1 {
		colorListText = args[1]
	}
	if !fixedFlag && checkRegexFlag {
		_, err := regexp.CompilePOSIX(pattern)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			os.Exit(1)
		}
		return
	}

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

	// NOTE: select submatch index func
	var findAllStringSubmatchIndex func(s string, n int) (indexes [][]int)
	if !fixedFlag {
		colorReg := regexp.MustCompilePOSIX(pattern)
		findAllStringSubmatchIndex = colorReg.FindAllStringSubmatchIndex
	} else {
		keyword := pattern
		// NOTE: for fixed string
		findAllStringSubmatchIndex = func(s string, n int) (indexes [][]int) {
			if keyword == "" {
				return nil
			}
			target := s
			start := 0
			for {
				ret := strings.Index(target, keyword)
				if ret >= 0 {
					start += ret
					end := start + len(keyword)
					// NOTE: for entire and 1st sub match
					indexes = append(indexes, []int{start, end, start, end})
					target = s[end:]
					start = end
					continue
				}
				break
			}
			return
		}
	}

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		text := scanner.Text()

		ansiText := cgrep.ANSITextParse(text)
		if verbose {
			fmt.Println("[START]")
			ansiText.Debug()
		}

		// NOTE: overwrite ansi color code
		m := findAllStringSubmatchIndex(ansiText.Plaintext, -1)
		if len(m) > 0 {
			for _, v := range m {
				// NOTE: v[0],v[1]: entire index set
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
					ansi := cgrep.ANSIEscapeCodeRange{
						Start: len([]rune(ansiText.Plaintext[:start])),
						End:   len([]rune(ansiText.Plaintext[:end])),
						Code:  string(code),
					}
					ansiText.ANSIRanges = append(ansiText.ANSIRanges, ansi)
				}
			}
		}

		// NOTE: only for debug
		if verbose {
			ansiText.DebugANSIRanges()
		}
		fmt.Println(ansiText)
		if verbose {
			fmt.Println("[END]")
		}
	}
	if err := scanner.Err(); err != nil {
		log.Printf("stdin read err:%s\n", err)
	}
}
