package js

import (
	"fmt"
	rs "github.com/pelmers/reshapify"
	"sort"
	"strings"
	"unicode"
)

type JSTokenizer struct {
	symbols  []rune
	keywords []string
	reserved []string
}

func NewJSTokenizer() *JSTokenizer {
	return &JSTokenizer{
		// symbols
		[]rune{'~', '?', ':', '!', '+', '-', '*', '/', '=', '(', ')', '[', ']', '{', '}', '<', '>', ';', ',', '.', '"', '\''},
		// keywords
		[]string{
			"break", "case", "comment", "continue", "default", "delete", "do", "else",
			"export", "false", "for", "function", "if", "import", "in", "instanceof", "label", "new", "null",
			"return", "switch", "this", "true", "typeof", "undefined", "var", "void", "while", "with",
		},
		// reserved
		[]string{
			"===", "==", "!==", "!=", "console", "log", "document", "window", "length", "push", "hasOwnProperty",
			">=", "<=", "+=", "-=", "reverse", "shift", "random", "Math", "floor",
		},
	}
}

func (tok *JSTokenizer) firstSymbol(text string) (int, rune) {
	// Return the index of the first symbol in text string
	for i, t := range text {
		if unicode.IsSpace(t) {
			return i, t
		}
		if rs.SearchRunes(tok.symbols, t) != -1 {
			return i, t
		}
	}
	// symbol not found, return special end code
	return len(text) - 1, rune(0)
}

// Return true if c is either a space or a symbol
func (tok *JSTokenizer) isSymbol(t string) bool {
	if len(t) != 1 {
		return false
	}
	r := rune(t[0])
	return unicode.IsSpace(r) || rs.SearchRunes(tok.symbols, r) != -1
}

// Return whether t is a javascript reserved word
func (tok *JSTokenizer) isKw(t string) bool {
	return rs.BSearchStrings(tok.keywords, strings.TrimSpace(t)) != -1
}

// Return whether the identifier t is reserved
func (tok *JSTokenizer) isReserved(t string) bool {
	return rs.BSearchStrings(tok.reserved, strings.TrimSpace(t)) != -1
}

// Return whether the string is quoted like a JS string
func (tok *JSTokenizer) isQuoted(a string) bool {
	return (strings.HasPrefix(a, "'") && strings.HasSuffix(a, "'")) || (strings.HasPrefix(a, "\"") && strings.HasSuffix(a, "\""))
}

// Arbitrarily rename non-reserved, non-keyword, non-symbol tokens in place
func (tok *JSTokenizer) RenameTokens(tokens []string, length int) {
	// first pass: create definition mapping
	defs := make(map[string]string)
	for _, t := range tokens {
		// rename if it is not a keyword, not a string, not reserved, and not a symbol
		if tok.isKw(t) || tok.isSymbol(t) || tok.isReserved(t) || rs.IsDigits(t) || tok.isQuoted(t) || len(t) == 0 {
			continue
		}
		// if not defined, define it
		if _, ok := defs[t]; !ok {
			// make sure the ident with unary operator attached doesn't change
			if strings.HasSuffix(t, "++") {
				if pre, ok := defs[strings.TrimSuffix(t, "++")]; ok {
					defs[t] = pre + "++"
				}
			} else if strings.HasSuffix(t, "--") {
				if pre, ok := defs[strings.TrimSuffix(t, "--")]; ok {
					defs[t] = pre + "--"
				}
			} else if strings.HasPrefix(t, "return ") {
				if suf, ok := defs[strings.TrimPrefix(t, "return ")]; ok {
					defs[t] = "return " + suf
				} else {
					defs[strings.TrimPrefix(t, "return ")] = rs.MakeIdentString(length)
					defs[t] = "return " + defs[suf]
				}
			} else {
				defs[t] = rs.MakeIdentString(length)
			}
		}
	}
	// second pass: replace tokens with their renames
	for i, _ := range tokens {
		if ren, ok := defs[tokens[i]]; ok {
			tokens[i] = ren
		}
	}
}

func (tok *JSTokenizer) Tokenize(code string, _reserved []string) []string {
	tok.reserved = append(tok.reserved, _reserved...)
	sort.Strings(tok.reserved)
	sort.Strings(tok.keywords)
	tokens := make([]string, 0)
	for {
		nsIndex, nsRune := tok.firstSymbol(code)
		if nsRune == rune(0) {
			break
		}
		tokens = append(tokens, code[:nsIndex])
		if !unicode.IsSpace(nsRune) {
			tokens = append(tokens, string(nsRune))
		}
		code = code[nsIndex+1:]
		// check for the speshul cases
		switch nsRune {
		// add the rest of the string as a token
		// TODO: support escaped quotes
		case '"', '\'':
			nsIndex = strings.IndexRune(code, nsRune)
			tokens[len(tokens)-1] = fmt.Sprintf("%c%s%c",
				nsRune, code[:nsIndex], nsRune)
			// advance the code to the close quote character
			code = code[nsIndex+1:]
		case '+', '-':
			// check for unary -- or ++
			if rune(code[0]) == nsRune {
				// move code forward 1 rune
				code = code[1:]
				// delete the - or + we just added
				tokens = tokens[:len(tokens)-1]
				// add -- or ++ to the end of the previous token
				tokens[len(tokens)-1] += string([]rune{nsRune, nsRune})
			}
			// check for += and -=
			if code[0] == '=' {
				tokens[len(tokens)-1] = string(nsRune) + "="
				code = code[1:]
			}
			// check for single-line comment
		case '/':
			if code[0] == '/' {
				// ignore the rest of the line and remove the / token we added
				tokens = tokens[:len(tokens)-1]
				code = code[strings.IndexRune(code, '\n')+1:]
			}
		case '=':
			if code[0] == '=' && code[1] == '=' {
				tokens[len(tokens)-1] = "==="
				code = code[2:]
			} else if code[0] == '=' {
				tokens[len(tokens)-1] = "=="
				code = code[1:]
			}
		case '!':
			if code[0] == '=' && code[1] == '=' {
				tokens[len(tokens)-1] = "!=="
				code = code[2:]
			} else if code[0] == '=' {
				tokens[len(tokens)-1] = "!="
				code = code[1:]
			}
		case '<', '>':
			if code[0] == '=' {
				tokens[len(tokens)-1] = string(nsRune) + "="
				code = code[1:]
			}
		}
		code = strings.TrimSpace(code)
	}
	// go through tokens and hardcode a space before and after any keywords
	for i, _ := range tokens[:len(tokens)-1] {
		if tok.isKw(tokens[i]) {
			if !tok.isSymbol(tokens[i+1]) {
				tokens[i] = tokens[i] + " "
				if tokens[i] == "return " {
					tokens[i] += tokens[i+1]
					tokens[i+1] = ""
				}
			}
			if i > 0 && !tok.isSymbol(tokens[i-1]) {
				tokens[i] = " " + tokens[i]
			}
		}
	}
	return tokens
}
