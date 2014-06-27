package js

import (
	"bufio"
	"fmt"
	"github.com/pelmers/cram/tokenize"
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

func TestTokenizeFile(t *testing.T) {
	filename := "test/test.js"
	file, err := os.Open(filename)
	if err != nil {
		t.Error("Could not open ", filename, " for testing")
	}
	reader := bufio.NewReader(file)
	code, err := ioutil.ReadAll(reader)
	if err != nil {
		t.Error("Could not convert code into string.")
	}
	tok := NewJSTokenizer()
	tokens := tok.Tokenize(string(code), []string{""})
	expected_tokens := []string{}
	if !tokenize.SlicesEqual(tokens, expected_tokens) {
		fmt.Println(strings.Join(tokens, ""))
		t.Error("Expected", expected_tokens, "Got", tokens)
	}
	file.Close()
}