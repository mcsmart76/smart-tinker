// Print a number as a big number.
//
// Inspired by the example in:
// http://www.drdobbs.com/open-source/getting-going-with-go/240004971

package main

import (
	"fmt"
	"log"
	"os"
)

// Each digit is hand-formatted.  I am not a professional typographer,
// but did my best to give each character a pleasant shape.  This is
// not a fixed-width font.
var bigNumbers = [][]string{
        {
                "  000  ",
                " 0   0 ",
                "0     0",
                "0     0",
                "0     0",
                " 0   0 ",
                "  000  "},
        {
                "  1 ",
                " 11 ",
                "  1 ",
                "  1 ",
                "  1 ",
                "  1 ",
                " 111"},
        {
                " 222 ",
                "2   2",
                "   2 ",
                "  2  ",
                " 2   ",
                "2    ",
                "22222"},
        {
                " 3333 ",
                "3    3",
                "     3",
                "  333 ",
                "     3",
                "3    3",
                " 3333 "},
        {
                "   44 ",
                "  4 4 ",
                " 4  4 ",
                "444444",
                "    4 ",
                "    4 ",
                "    4 "},
        {
                "555555",
                "5     ",
                "5     ",
                " 5555 ",
                "     5",
                "5    5",
                " 5555 "},
        {
                " 6666 ",
                "6    6",
                "6     ",
                "66666 ",
                "6    6",
                "6    6",
                " 6666 "},
	{
                "777777",
                "     7",
                "    7 ",
                "   7  ",
                "  7   ",
                " 7    ",
                "7     "},
        {
                " 8888 ",
                "8    8",
                "8    8",
                " 8888 ",
                "8    8",
                "8    8",
                " 8888 "},
        {
                " 9999 ",
                "9    9",
                "9    9",
                " 99999",
                "     9",
                "     9",
                "     9"},
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: bignumber <number>")
		os.Exit(1)
	}

	// Verify each character has the same height and each
	// character has an equal width per row.
	verified := true
	height := len(bigNumbers[0])
	for d, digit := range bigNumbers {
		if len(digit) != height {
			log.Printf("The character set is invalid; " +
				   "%v is a different height", d)
			verified = false
		}

		width := len(digit[0])
		for _, row := range digit {
			if len(row) != width {
				log.Printf("The character set is invalid; " +
					   "%v has rows of different widths", d)
				verified = false
				break
			}
		}
	}
	if !verified {
		// Printing out all the errors lets the programmer fix
		// the "font" quicker.
		os.Exit(1)
	}

	numberAsString := os.Args[1]
	for row := range bigNumbers[0] {
		line := ""
		for i := range numberAsString {
			character := numberAsString[i] - '0'
			if character < 0 || character > 9 {
				log.Fatal("Invalid number; try 0-9 or send a patch")
			}
			line += bigNumbers[character][row] + " "
		}
		fmt.Println(line)
	}
}

