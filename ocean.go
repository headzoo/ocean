/*
Ocean
=====
A simple lexer for go that supports shell-style quoting, commenting, piping,
redirecting, and escaping.

Code originally forked from go-shlex (http://code.google.com/p/go-shlex/).
Contributions made by flynn/go-shlex (https://github.com/flynn/go-shlex).


License
=======
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
	"io"
	"strings"
)

// TokenType is a top-level token; a word, space, comment, unknown.
type TokenType string

// TokenValue is the value of the token, usually a string.
type TokenValue string

// RuneType is the type of a UTF-8 character; a character, quote, space, escape.
type RuneType string

// RuneTypeMap is a map of RuneTokeType values.
type RuneTypeMap map[rune]RuneType

// LexerState is used within the lexer state machine to keep track of the current state.
type LexerState int

// Token represents a single "token" found within a stream.
type Token struct {
	tokenType TokenType
	value     TokenValue
}

const (
	CLASS_CHAR              = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789#._-,/@$*()+=:;&^%~"
	CLASS_SPACE             = " \t\r\n"
	CLASS_ESCAPING_QUOTE    = "\""
	CLASS_NONESCAPING_QUOTE = "'"
	CLASS_ESCAPE            = "\\"
	CLASS_PIPE              = "|"
	CLASS_REDIRECT          = "><"

	RUNE_UNKNOWN      RuneType = "UNKNOWN"
	RUNE_CHAR         RuneType = "CHAR"
	RUNE_SPACE        RuneType = "SPACE"
	RUNE_QUOTE_DOUBLE RuneType = "QUOTE_DOUBLE"
	RUNE_QUOTE_SINGLE RuneType = "QUOTE_SINGLE"
	RUNE_ESCAPE       RuneType = "ESCAPE"
	RUNE_PIPE         RuneType = "PIPE"
	RUNE_REDIRECT     RuneType = "REDIRECT"
	RUNE_EOF          RuneType = "EOF"

	TOKEN_UNKNOWN  TokenType = "UNKNOWN"
	TOKEN_WORD     TokenType = "WORD"
	TOKEN_SPACE    TokenType = "SPACE"
	TOKEN_PIPE     TokenType = "PIPE"
	TOKEN_REDIRECT TokenType = "REDIRECT"

	INITIAL_TOKEN_CAPACITY = 100
)

// Split splits a string in to a slice of strings, based upon shell-style rules for
// quoting, escaping, and spaces.
func Split(s string) ([]TokenValue, error) {
	l, err := NewLexer(strings.NewReader(s))
	if err != nil {
		return nil, err
	}
	subStrings := []TokenValue{}
	for {
		word, err := l.NextWord()
		if err != nil {
			if err == io.EOF {
				return subStrings, nil
			}
			return subStrings, err
		}
		subStrings = append(subStrings, word)
	}

	return subStrings, nil
}
