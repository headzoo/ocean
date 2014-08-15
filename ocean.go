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
type TokenType int

// TokenValue is the value of the token, usually a string.
type TokenValue string

// RuneType is the type of a UTF-8 character; a character, quote, space, escape.
type RuneType int

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
	CLASS_CHAR              = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789._-,/@$*()+=:;&^%~"
	CLASS_SPACE             = " \t\r\n"
	CLASS_ESCAPING_QUOTE    = "\""
	CLASS_NONESCAPING_QUOTE = "'"
	CLASS_ESCAPE            = "\\"
	CLASS_COMMENT           = "#"
	CLASS_PIPE              = "|"
	CLASS_REDIRECT          = "><"

	RUNE_UNKNOWN      RuneType = 0
	RUNE_CHAR         RuneType = 1
	RUNE_SPACE        RuneType = 2
	RUNE_QUOTE_DOUBLE RuneType = 3
	RUNE_QUOTE_SINGLE RuneType = 4
	RUNE_ESCAPE       RuneType = 5
	RUNE_COMMENT      RuneType = 6
	RUNE_PIPE         RuneType = 7
	RUNE_REDIRECT     RuneType = 8
	RUNE_EOF          RuneType = 9

	TOKEN_UNKNOWN  TokenType = 0
	TOKEN_WORD     TokenType = 1
	TOKEN_SPACE    TokenType = 2
	TOKEN_COMMENT  TokenType = 3
	TOKEN_PIPE     TokenType = 4
	TOKEN_REDIRECT TokenType = 5

	STATE_START           LexerState = 0
	STATE_APPEND          LexerState = 1
	STATE_ESCAPING        LexerState = 2
	STATE_ESCAPING_QUOTED LexerState = 3
	STATE_QUOTED_ESCAPING LexerState = 4
	STATE_QUOTED          LexerState = 5
	STATE_COMMENT         LexerState = 6
	STATE_EMIT            LexerState = 7

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
