/*
Copyright 2012 Google Inc. All Rights Reserved.

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

func checkError(err error, t *testing.T) {
	if err != nil {
		t.Error(err)
	}
}

func TestClassifier(t *testing.T) {
	classifier := NewDefaultClassifier()
	runeTests := map[int32]RuneTokenType{
		'a':  RUNETOKEN_CHAR,
		' ':  RUNETOKEN_SPACE,
		'"':  RUNETOKEN_ESCAPING_QUOTE,
		'\'': RUNETOKEN_NONESCAPING_QUOTE,
		'#':  RUNETOKEN_COMMENT,
		'|': RUNETOKEN_PIPE,
		'>': RUNETOKEN_REDIRECT_OUT,
		'<': RUNETOKEN_REDIRECT_IN,
	}
	for rune, expectedType := range runeTests {
		foundType := classifier.ClassifyRune(rune)
		if foundType != expectedType {
			t.Logf("Expected type: %v for rune '%c'(%v). Found type: %v.", expectedType, rune, rune, foundType)
			t.Fail()
		}
	}
}

func TestTokenizer(t *testing.T) {
	testInput := strings.NewReader("one two \"three four\" \"five \\\"six\\\"\" seven#eight # nine # ten\n eleven | sixteen > \"seventeen eighteen\" < nineteen")
	expectedTokens := []*Token{
		&Token{
			tokenType: TOKEN_WORD,
			value:     "one"},
		&Token{
			tokenType: TOKEN_WORD,
			value:     "two"},
		&Token{
			tokenType: TOKEN_WORD,
			value:     "three four"},
		&Token{
			tokenType: TOKEN_WORD,
			value:     "five \"six\""},
		&Token{
			tokenType: TOKEN_WORD,
			value:     "seven#eight"},
		&Token{
			tokenType: TOKEN_COMMENT,
			value:     " nine # ten"},
		&Token{
			tokenType: TOKEN_WORD,
			value:     "eleven"},
		&Token{
			tokenType: TOKEN_PIPE,
			value:     "|"},
		&Token{
			tokenType: TOKEN_WORD,
			value:     "sixteen"},
		&Token{
			tokenType: TOKEN_REDIRECT_OUT,
			value:	   ">"},
		&Token{
			tokenType: TOKEN_WORD,
			value:     "seventeen eighteen"},
		&Token{
			tokenType: TOKEN_REDIRECT_IN,
			value:     "<"},
		&Token{
			tokenType: TOKEN_WORD,
			value:     "nineteen"},
		}

	tokenizer, err := NewTokenizer(testInput)
	checkError(err, t)
	for _, expectedToken := range expectedTokens {
		foundToken, err := tokenizer.NextToken()
		checkError(err, t)
		if !foundToken.Equal(expectedToken) {
			t.Error("Expected token:", expectedToken, ". Found:", foundToken)
		}
	}
}

func TestLexer(t *testing.T) {
	testInput := strings.NewReader("one")
	expectedWord := "one"
	lexer, err := NewLexer(testInput)
	checkError(err, t)
	foundWord, err := lexer.NextWord()
	checkError(err, t)
	if expectedWord != foundWord {
		t.Error("Expected word:", expectedWord, ". Found:", foundWord)
	}
}

func TestSplitSimple(t *testing.T) {
	assertSplit(
		"one two three",
		[]string{"one", "two", "three"},
		t,
	)
}

func TestSplitEscapingQuotes(t *testing.T) {
	assertSplit(
		"one \"two three\" four",
		[]string{"one", "two three", "four"},
		t,
	)
}

func TestSplitNonEscapingQuotes(t *testing.T) {
	assertSplit(
		"one 'two three' four",
		[]string{"one", "two three", "four"},
		t,
	)
}

func TestSplitPiped(t *testing.T) {
	assertSplit(
		"one two|three four",
		[]string{"one", "two", "|", "three", "four"},
		t,
	)
}

func TestSplitRedirectOut(t *testing.T) {
	assertSplit(
		"one two > three.txt",
		[]string{"one", "two", ">", "three.txt"},
		t,
	)
}

func TestSplitRedirectIn(t *testing.T) {
	assertSplit(
		"one < two.txt",
		[]string{"one", "<", "two.txt"},
		t,
	)
}
/*
func TestSplitComplex(t *testing.T) {
	assertSplit(
		"ls -l /|grep 'foo.txt' >> saved.txt",
		[]string{"ls", "-l", "/", "|", "grep", "foo.txt", ">>", "saved.txt"},
		t,
	)
}
*/


// assertSplit splits the string, and asserts it's equal to the expected value.
func assertSplit(input string, expected []string, t *testing.T) {
	actual, err := Split(input)
	if err != nil {
		t.Error("Split returned error:", err)
	}
	assertSliceEquals(expected, actual, t)
}

// assertSliceEquals asserts that two slices are equal.
func assertSliceEquals(expected, actual []string, t *testing.T) {
	if len(expected) != len(actual) {
		t.Error("Split expected:", len(expected), "results. Found:", len(actual), "results")
	}
	for i := range actual {
		if actual[i] != expected[i] {
			t.Error("Item:", i, "(", actual[i], ") differs from the expected value:", expected[i])
		}
	}
}
