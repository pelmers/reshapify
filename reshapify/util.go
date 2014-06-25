package main

import (
	"fmt"
	"log"
	"math/rand"
	"strings"
)

// Return the sum of lengths of all strings in a
func TotalLength(a []string) int {
	s := 0
	for _, n := range a {
		s += len(n)
	}
	return s
}

// Given a string of tokens and a function which maps
// line number -> desired width,
// split tokens into lines which do not exceed desired width,
// unless the desired width is less than the length of one token.
func SplitLines(tokens []string, widthFromLineNo func(int) int) [][]string {
	lines := make([][]string, 0)
	line_no := 0
	token_no := 0
	for token_no < len(tokens) {
		lines = append(lines, make([]string, 0))
		width := widthFromLineNo(line_no)
		if width <= 0 {
			log.Printf("Negative width, defaulting to 1 : %d on line %d\n", width, line_no)
			width = 1
		}
		for TotalLength(lines[line_no]) < width {
			lines[line_no] = append(lines[line_no], tokens[token_no])
			token_no++
			if token_no == len(tokens) {
				return lines
			}
		}
		// advance line number and take off the last token of previous line
		// since the last token pushed the string over the square width
		// unless the last line was only one token long
		if len(lines[line_no]) > 1 {
			lines[line_no] = lines[line_no][:len(lines[line_no])-1]
			token_no--
		}
		line_no++
	}
	return lines
}

// Given a slice of lines, where each line is a slice of token strings that should
// appear on that line and a function that maps line number -> desired width,
// add spaces to each line to make it reach the desired width if possible.
// If centered is set to true, also center the output.
// Join the justified lines together and return a string.
func JustifyByWidth(lines [][]string, widthFromLineNo func(int) int, centered bool) string {
	var maxWidth int
	if centered {
		for _, line := range lines {
			if len(line) > maxWidth {
				maxWidth = len(line)
			}
		}
	}
	justifiedLines := make([]string, 0, len(lines))
	for i, line := range lines {
		width := widthFromLineNo(i)
		if width <= 0 {
			log.Printf("Negative width, defaulting to 1 : %d on line %d\n", width, i)
			width = 1
		}
		for TotalLength(line) < width {
			idx := rand.Intn(len(line))
			line[idx] += " "
		}
		spacing := ""
		// center by prepending spaces such that the center is at maxWidth/2
		if centered {
			spacing = fmt.Sprintf("%*s", (maxWidth - width/2), " ")
		}
		justifiedLines = append(justifiedLines, spacing+strings.Join(line, ""))
	}
	return strings.Join(justifiedLines, "\n")
}
