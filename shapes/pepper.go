package shapes

import (
	"fmt"
	"math/rand"
	"strings"
)

// Return the number of tokens it takes to exceed the totalLength
func countTokensUpToTotalLength(tokens []string, totalLength int) int {
	var sofar int
	for i, t := range tokens {
		sofar += len(t)
		if sofar > totalLength {
			return i
		}
	}
	return len(tokens)
}

// Given a string, randomly adjust indendation by up to ratio amount.
func perturbIndents(text string, round float64) string {
	textSplit := strings.Split(text, "\n")
	for i, line := range textSplit {
		line = strings.TrimLeft(line, " ")
		spaces := len(textSplit[i]) - len(line)
		// random number in range (1-round, 1+round)
		perturb := (1.0 - round) + rand.Float64()*round
		textSplit[i] = fmt.Sprintf("%*s%s", int(float64(spaces)*perturb), " ", line)
	}
	return strings.Join(textSplit, "\n")
}

func Pepper(tokens []string, ratio float64) string {
	// model: http://upload.wikimedia.org/wikipedia/commons/thumb/8/86/Habanero_closeup_edit2.jpg/640px-Habanero_closeup_edit2.jpg
	area := float64(TotalLength(tokens))
	// How to apportion the area:
	stem := 0.03
	square := 0.5
	// rest: flipped triangle
	stemL := countTokensUpToTotalLength(tokens, int(area*stem))
	stemTokens := tokens[:stemL]
	squareL := countTokensUpToTotalLength(tokens[stemL:], int(area*square))
	squareTokens := tokens[stemL : stemL+squareL]
	triTokens := tokens[stemL+squareL:]
	// the top part is a very skinny curved stem
	squareCode := Square(squareTokens, ratio)
	// the width of the square serves as the base of the triangle
	base := float64(strings.Index(squareCode, "\n"))
	triArea := float64(TotalLength(triTokens))
	triParams := parametrizedTriangle(triArea, base, 2.0*triArea/base, 5)
	triWidth := func(h int) int {
		return triParams(int(2.0*triArea/base - float64(h)))
	}
	triCode := JustifyByWidth(SplitLines(triTokens, triWidth), triWidth, true)
	// (try to) gently curve the stem into the middle of the base
	offset := int(base * 0.15) // don't start the stem right on the edge
	slope := (base/2 - float64(offset)) / float64(len(stemTokens))
	stemCode := ""
	for i, _ := range stemTokens[:len(stemTokens)-1] {
		if i%2 == 0 {
			stemCode += fmt.Sprintf("%*s%s%s\n", offset+int(slope*float64(i)), " ",
				stemTokens[i], stemTokens[i+1])
		}
	}
	return strings.Join([]string{
		perturbIndents(stemCode, 0.3) + squareCode,
		perturbIndents(triCode, 0.5),
	}, "\n")
}