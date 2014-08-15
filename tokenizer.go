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
	"bufio"
	"bytes"
	"errors"
	
	"io"
)

// TokenClass is a top-level token; a word, space, comment, unknown.
type TokenClass string

// TokenValue is the value of the token, usually a string.
type TokenValue string

// TokenState is used within the lexer state machine to keep track of the current state.
type TokenState int

const (
	TOKEN_UNKNOWN  TokenClass = "UNKNOWN"
	TOKEN_WORD     TokenClass = "WORD"
	TOKEN_SPACE    TokenClass = "SPACE"
	TOKEN_PIPE     TokenClass = "PIPE"
	TOKEN_REDIRECT TokenClass = "REDIRECT"

	STATE_START           TokenState = 0
	STATE_WORD            TokenState = 1
	STATE_ESCAPING        TokenState = 2
	STATE_ESCAPING_QUOTED TokenState = 3
	STATE_QUOTED_ESCAPING TokenState = 4
	STATE_QUOTED          TokenState = 5
	STATE_COMMENT         TokenState = 6
	STATE_EMIT            TokenState = 7
)

// Token represents a single "token" found within a stream.
type Token struct {
	Class TokenClass
	Value TokenValue
}

// Creates and returns a new Token instance.
func NewToken(class TokenClass, value TokenValue) *Token {
	return &Token{
		Class: class,
		Value: value,
	}
}

// Tokenizer turns an input stream in to a sequence of typed tokens.
type Tokenizer struct {
	input      *bufio.Reader
	classifier *Classifier
}

// Creates and returns a new tokenizer.
func NewTokenizer(reader io.Reader) *Tokenizer {
	return &Tokenizer{
		input:      bufio.NewReader(reader),
		classifier: NewClassifier(),
	}
}

// NextToken returns the next token in the stream, and an error value. If there are no more
// tokens available, the error value will be io.EOF.
func (tokenizer *Tokenizer) NextToken() (*Token, error) {
	var (
		err       error
		next      rune
		class     RuneClass
		tokenType TokenClass
	)

	state := STATE_START
	buffer := bytes.NewBuffer(make([]byte, 0, 1000))

SCAN:
	for state != STATE_EMIT {
		next, class, err = tokenizer.readRune()
		if err != nil {
			return nil, err
		}

		switch state {
		case STATE_START: // no runes read yet
			{
				switch class {
				case RUNE_EOF:
					{
						return nil, io.EOF
					}
				case RUNE_CHAR:
					{
						tokenType = TOKEN_WORD
						buffer.WriteRune(next)
						state = STATE_WORD
					}
				case RUNE_SPACE:
					{
					}
				case RUNE_QUOTE_DOUBLE:
					{
						tokenType = TOKEN_WORD
						state = STATE_QUOTED_ESCAPING
					}
				case RUNE_QUOTE_SINGLE:
					{
						tokenType = TOKEN_WORD
						state = STATE_QUOTED
					}
				case RUNE_ESCAPE:
					{
						tokenType = TOKEN_WORD
						state = STATE_ESCAPING
					}
				case RUNE_PIPE:
					{
						tokenType = TOKEN_PIPE
						buffer.WriteRune(next)
						state = STATE_EMIT
					}
				case RUNE_REDIRECT:
					{
						tokenType = TOKEN_REDIRECT
						buffer.WriteRune(next)
						state = STATE_EMIT

						n, _, e := tokenizer.readRune()
						if e != nil {
							return nil, err
						}
						if n == next {
							buffer.WriteRune(next)
						} else {
							tokenizer.unreadRune()
						}
					}
				default:
					{
						return nil, errorf("Unknown rune: %v", next)
					}
				}
			}
		case STATE_WORD: // in a regular word
			{
				switch class {
				case RUNE_EOF:
					{
						break SCAN
					}
				case RUNE_CHAR:
					{
						buffer.WriteRune(next)
					}
				case RUNE_SPACE:
					{
						tokenizer.unreadRune()
						break SCAN
					}
				case RUNE_QUOTE_DOUBLE:
					{
						state = STATE_QUOTED_ESCAPING
					}
				case RUNE_QUOTE_SINGLE:
					{
						state = STATE_QUOTED
					}
				case RUNE_ESCAPE:
					{
						state = STATE_ESCAPING
					}
				case RUNE_PIPE:
					{
						tokenizer.unreadRune()
						state = STATE_EMIT
					}
				case RUNE_REDIRECT:
					{
						tokenizer.unreadRune()
						state = STATE_EMIT
					}
				default:
					{
						return nil, errorf("Uknown rune: %v", next)
					}
				}
			}
		case STATE_ESCAPING: // the next rune after an escape character
			{
				switch class {
				case RUNE_EOF:
					{
						err = errors.New("EOF found after escape character")
						break SCAN
					}
				case RUNE_CHAR, RUNE_SPACE, RUNE_QUOTE_DOUBLE, RUNE_QUOTE_SINGLE, RUNE_ESCAPE,
					RUNE_PIPE, RUNE_REDIRECT:
					{
						state = STATE_WORD
						buffer.WriteRune(next)
					}
				default:
					{
						return nil, errorf("Uknown rune: %v", next)
					}
				}
			}
		case STATE_ESCAPING_QUOTED: // the next rune after an escape character, in double quotes
			{
				switch class {
				case RUNE_EOF:
					{
						err = errors.New("EOF found after escape character")
						break SCAN
					}
				case RUNE_CHAR, RUNE_SPACE, RUNE_QUOTE_DOUBLE, RUNE_QUOTE_SINGLE, RUNE_ESCAPE,
					RUNE_PIPE, RUNE_REDIRECT:
					{
						state = STATE_QUOTED_ESCAPING
						buffer.WriteRune(next)
					}
				default:
					{
						return nil, errorf("Uknown rune: %v", next)
					}
				}
			}
		case STATE_QUOTED_ESCAPING: // in escaping double quotes
			{
				switch class {
				case RUNE_EOF:
					{
						err = errors.New("EOF found when expecting closing quote.")
						break SCAN
					}
				case RUNE_CHAR, RUNE_SPACE, RUNE_QUOTE_SINGLE,
					RUNE_PIPE, RUNE_REDIRECT:
					{
						buffer.WriteRune(next)
					}
				case RUNE_QUOTE_DOUBLE:
					{
						state = STATE_WORD
					}
				case RUNE_ESCAPE:
					{
						state = STATE_ESCAPING_QUOTED
					}
				default:
					{
						return nil, errorf("Uknown rune: %v", next)
					}
				}
			}
		case STATE_QUOTED: // in non-escaping single quotes
			{
				switch class {
				case RUNE_EOF:
					{
						err = errors.New("EOF found when expecting closing quote.")
						break SCAN
					}
				case RUNE_CHAR, RUNE_SPACE, RUNE_QUOTE_DOUBLE, RUNE_ESCAPE,
					RUNE_PIPE, RUNE_REDIRECT:
					{
						buffer.WriteRune(next)
					}
				case RUNE_QUOTE_SINGLE:
					{
						state = STATE_WORD
					}
				default:
					{
						return nil, errorf("Uknown rune: %v", next)
					}
				}
			}
		case STATE_COMMENT:
			{
				switch class {
				case RUNE_EOF:
					{
						break SCAN
					}
				case RUNE_CHAR, RUNE_QUOTE_DOUBLE, RUNE_ESCAPE,
					RUNE_QUOTE_SINGLE, RUNE_PIPE, RUNE_REDIRECT:
					{
						buffer.WriteRune(next)
					}
				case RUNE_SPACE:
					{
						if next == '\n' {
							state = STATE_START
							break SCAN
						} else {
							buffer.WriteRune(next)
						}
					}
				default:
					{
						return nil, errorf("Uknown rune: %v", next)
					}
				}
			}
		default:
			{
				panic(errorf("Unexpected state: %v", state))
			}
		}
	}

	return NewToken(tokenType, TokenValue(buffer.String())), err
}

// readRune returns the next rune in the stream.
func (tokenizer *Tokenizer) readRune() (rune, RuneClass, error) {
	next, _, err := tokenizer.input.ReadRune()
	class := tokenizer.classifier.Classify(next)

	if err != nil {
		if err == io.EOF {
			class = RUNE_EOF
			err = nil
		}
	}

	return next, class, err
}

// unread puts the previously read rune back into the stream.
func (tokenizer *Tokenizer) unreadRune() {
	tokenizer.input.UnreadRune()
}
