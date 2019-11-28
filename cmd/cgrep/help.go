package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/mgutz/ansi"
)

func help() {
	fmt.Fprintf(os.Stderr, ""+
		os.Args[0]+` <regex pattern with ()> [colors with ,]
  e.g.
    echo 'float x = 10.0f;' | cgrep '(([0-9]+)\.([0-9]+)f)' 'magenta'
    echo 'g++ -std=c++11 -o main.o -c main.cpp' | cgrep '(-o [^ ]+)' 'yellow+b'

`)
	flag.PrintDefaults()
	fmt.Fprintln(os.Stderr)
	fmt.Fprintf(os.Stderr, `
## color format
`[1:])
	// NOTE: over kill printing
	// ansi.PrintStyles()
	printColorSamples()
	fmt.Fprintf(os.Stderr, `
FYI: https://github.com/mgutz/ansi
`[1:])
}

func printColorSamples() {
	colors := []string{
		"reset",
		"off",
		"black",
		"red",
		"green",
		"yellow",
		"blue",
		"magenta",
		"cyan",
		"white",
	}
	fmt.Fprintf(os.Stderr, "* colors      :[")
	for i, color := range colors {
		if i > 0 {
			fmt.Fprintf(os.Stderr, ", ")
		}
		fmt.Fprintf(os.Stderr, "%s", ansi.Color(color, color))
	}
	fmt.Fprintln(os.Stderr, "]")

	fmt.Fprintf(os.Stderr, "* with options:[")
	fmt.Fprintf(os.Stderr, "%s %s", ansi.Color("red", "red"), ansi.Reset)                     // red
	fmt.Fprintf(os.Stderr, "%s %s", ansi.Color("red+b", "red+b"), ansi.Reset)                 // red bold
	fmt.Fprintf(os.Stderr, "%s %s", ansi.Color("red+B", "red+B"), ansi.Reset)                 // red blinking
	fmt.Fprintf(os.Stderr, "%s %s", ansi.Color("red+u", "red+u"), ansi.Reset)                 // red underline
	fmt.Fprintf(os.Stderr, "%s %s", ansi.Color("red+bh", "red+bh"), ansi.Reset)               // red bold bright
	fmt.Fprintf(os.Stderr, "%s %s", ansi.Color("red:white", "red:white"), ansi.Reset)         // red on white
	fmt.Fprintf(os.Stderr, "%s %s", ansi.Color("red+b:white+h", "red+b:white+h"), ansi.Reset) // red bold on white bright
	fmt.Fprintf(os.Stderr, "%s %s", ansi.Color("red+B:white+h", "red+B:white+h"), ansi.Reset) // red blink on white bright
	fmt.Fprintln(os.Stderr, "]")
}
