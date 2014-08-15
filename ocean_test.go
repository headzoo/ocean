/*
Copyright 2014 Sean Hickey <sean@headzoo.io>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package ocean

import (
	"strings"
	"testing"
)

func TestClassifier(test *testing.T) {
	classifier := NewClassifier()
	runeTests := RuneTypeMap{
		'a':  RUNE_CHAR,
		' ':  RUNE_SPACE,
		'"':  RUNE_QUOTE_DOUBLE,
		'\'': RUNE_QUOTE_SINGLE,
		'|':  RUNE_PIPE,
		'>':  RUNE_REDIRECT,
	}

	for rune, expected := range runeTests {
		actual := classifier.Classify(rune)
		if actual != expected {
			test.Error("Expected type: %v for rune '%c'(%v). Found type: %v.", expected, rune, rune, actual)
		}
	}
}

func TestTokenizer(test *testing.T) {
	input := strings.NewReader(`one two "three four" "five \"six\"" \n eleven | sixteen > "seventeen eighteen" < nineteen`)
	expected := []*Token{
		NewToken(TOKEN_WORD, "one"),
		NewToken(TOKEN_WORD, "two"),
		NewToken(TOKEN_WORD, "three four"),
		NewToken(TOKEN_WORD, `five \"six\"`),
		NewToken(TOKEN_WORD, "eleven"),
		NewToken(TOKEN_PIPE, "|"),
		NewToken(TOKEN_WORD, "sixteen"),
		NewToken(TOKEN_REDIRECT, ">"),
		NewToken(TOKEN_WORD, "seventeen eighteen"),
		NewToken(TOKEN_REDIRECT, "<"),
		NewToken(TOKEN_WORD, "nineteen"),
	}

	tokenizer := NewTokenizer(input)
	for _, ex := range expected {
		actual, err := tokenizer.NextToken()
		assertNilError(err, test)
		if !actual.Equal(actual) {
			test.Error("Expected: ", ex, ", Actual: ", actual)
		}
	}
}

func TestLexer(test *testing.T) {
	input := strings.NewReader("one")
	expected := TokenValue("one")

	actual, err := NewLexer(input).NextWord()
	assertNilError(err, test)

	if expected != actual {
		test.Error("Expected word:", expected, ". Found:", actual)
	}
}

func TestSplitSimple(test *testing.T) {
	assertSplit(
		`one two three`,
		[]TokenValue{"one", "two", "three"},
		test,
	)
}

func TestSplitEscapingQuotes(test *testing.T) {
	assertSplit(
		`one "two three" four`,
		[]TokenValue{"one", "two three", "four"},
		test,
	)
}

func TestSplitNonEscapingQuotes(test *testing.T) {
	assertSplit(
		`one 'two three' four`,
		[]TokenValue{"one", "two three", "four"},
		test,
	)
}

func TestSplitPiped(test *testing.T) {
	assertSplit(
		`one two|three four`,
		[]TokenValue{"one", "two", "|", "three", "four"},
		test,
	)
}

func TestSplitRedirectOut(test *testing.T) {
	assertSplit(
		`one two > three.txt`,
		[]TokenValue{"one", "two", ">", "three.txt"},
		test,
	)
}

func TestSplitRedirectIn(test *testing.T) {
	assertSplit(
	`one < two.txt`,
	[]TokenValue{"one", "<", "two.txt"},
	test,
)
}

func TestSplitComplex(test *testing.T) {
assertSplit(
	`ls -l /|find << "find.txt" | grep 'foo.txt' >> saved.txt`,
		[]TokenValue{"ls", "-l", "/", "|", "find", "<<", "find.txt", "|", "grep", "foo.txt", ">>", "saved.txt"},
		test,
	)
}

// assertSplit splits the string, and asserts it's equal to the expected value.
func assertSplit(input string, expected []TokenValue, test *testing.T) {
	actual, err := Split(input)
	if err != nil {
		test.Error("Split returned error:", err)
	}
	assertTokenValueEquals(expected, actual, test)
}

// assertSliceEquals asserts that two slices are equal.
func assertTokenValueEquals(expected, actual []TokenValue, test *testing.T) {
	if len(expected) != len(actual) {
		test.Error("Split expected:", len(expected), "results. Found:", len(actual), "results")
		return
	}
	for i := range actual {
		if actual[i] != expected[i] {
			test.Error("Item:", i, "(", actual[i], ") differs from the expected value:", expected[i])
		}
	}
}

func assertNilError(err error, test *testing.T) {
	if err != nil {
		test.Error(err)
	}
}

// Equal returns a boolean value indicating whether two tokens are equal.
// Two tokens are equal if both their types and values are equal. A nil token can
// never equal another token.
func (a *Token) Equal(b *Token) bool {
	if a == nil || b == nil {
		return false
	}
	if a.Class != b.Class {
		return false
	}

	return a.Value == b.Value
}
