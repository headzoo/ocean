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
	"errors"
	"fmt"
	"io"
)

const (
	STATE_START           LexerState = 0
	STATE_APPEND          LexerState = 1
	STATE_ESCAPING        LexerState = 2
	STATE_ESCAPING_QUOTED LexerState = 3
	STATE_QUOTED_ESCAPING LexerState = 4
	STATE_QUOTED          LexerState = 5
	STATE_COMMENT         LexerState = 6
	STATE_EMIT            LexerState = 7
)

// Tokenizer turns an input stream in to a sequence of typed tokens.
type Tokenizer struct {
	input      *bufio.Reader
	classifier *Classifier
}

// Creates and returns a new tokenizer.
func NewTokenizer(reader io.Reader) (*Tokenizer, error) {
	input := bufio.NewReader(reader)
	classifier := NewClassifier()
	tokenizer := &Tokenizer{
		input:      input,
		classifier: classifier,
	}

	return tokenizer, nil
}

// scanStream scans the stream for the next token.
// This uses an internal state machine. It will panic if it encounters a character
// which it does not know how to handle.
func (t *Tokenizer) scanStream() (*Token, error) {
	var (
		tokenType    TokenType
		nextRune     rune
		nextRuneType RuneType
		err          error
	)

	state := STATE_START
	value := make([]rune, 0, INITIAL_TOKEN_CAPACITY)

SCAN:
	for state != STATE_EMIT {
		nextRune, _, err = t.input.ReadRune()
		nextRuneType = t.classifier.Classify(nextRune)
		if err != nil {
			if err == io.EOF {
				nextRuneType = RUNE_EOF
				err = nil
			} else {
				return nil, err
			}
		}

		switch state {
		case STATE_START: // no runes read yet
			{
				switch nextRuneType {
				case RUNE_EOF:
					{
						return nil, io.EOF
					}
				case RUNE_CHAR:
					{
						tokenType = TOKEN_WORD
						value = append(value, nextRune)
						state = STATE_APPEND
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
						value = append(value, nextRune)						
						state = STATE_EMIT
					}
				case RUNE_REDIRECT:
					{
						tokenType = TOKEN_REDIRECT
						value = append(value, nextRune)
						nr, _, err := t.input.ReadRune()
						if err != nil {
							if err == io.EOF {
								nextRuneType = RUNE_EOF
								err = nil
							} else {
								return nil, err
							}
						}
						if nr == nextRune {
							value = append(value, nr)
						} else {
							t.input.UnreadRune()
						}
						state = STATE_EMIT
					}
				default:
					{
						return nil, errors.New(fmt.Sprintf("Unknown rune: %v", nextRune))
					}
				}
			}
		case STATE_APPEND: // in a regular word
			{
				switch nextRuneType {
				case RUNE_EOF:
					{
						break SCAN
					}
				case RUNE_CHAR:
					{
						value = append(value, nextRune)
					}
				case RUNE_SPACE:
					{
						t.input.UnreadRune()
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
						t.input.UnreadRune()
						state = STATE_EMIT
					}
				case RUNE_REDIRECT:
					{
						t.input.UnreadRune()
						state = STATE_EMIT
					}
				default:
					{
						return nil, errors.New(fmt.Sprintf("Uknown rune: %v", nextRune))
					}
				}
			}
		case STATE_ESCAPING: // the next rune after an escape character
			{
				switch nextRuneType {
				case RUNE_EOF:
					{
						err = errors.New("EOF found after escape character")
						break SCAN
					}
				case RUNE_CHAR, RUNE_SPACE, RUNE_QUOTE_DOUBLE, RUNE_QUOTE_SINGLE, RUNE_ESCAPE,
					RUNE_PIPE, RUNE_REDIRECT:
					{
						state = STATE_APPEND
						value = append(value, nextRune)
					}
				default:
					{
						return nil, errors.New(fmt.Sprintf("Uknown rune: %v", nextRune))
					}
				}
			}
		case STATE_ESCAPING_QUOTED: // the next rune after an escape character, in double quotes
			{
				switch nextRuneType {
				case RUNE_EOF:
					{
						err = errors.New("EOF found after escape character")
						break SCAN
					}
				case RUNE_CHAR, RUNE_SPACE, RUNE_QUOTE_DOUBLE, RUNE_QUOTE_SINGLE, RUNE_ESCAPE,
					RUNE_PIPE, RUNE_REDIRECT:
					{
						state = STATE_QUOTED_ESCAPING
						value = append(value, nextRune)
					}
				default:
					{
						return nil, errors.New(fmt.Sprintf("Uknown rune: %v", nextRune))
					}
				}
			}
		case STATE_QUOTED_ESCAPING: // in escaping double quotes
			{
				switch nextRuneType {
				case RUNE_EOF:
					{
						err = errors.New("EOF found when expecting closing quote.")
						break SCAN
					}
				case RUNE_CHAR, RUNE_SPACE, RUNE_QUOTE_SINGLE,
					RUNE_PIPE, RUNE_REDIRECT:
					{
						value = append(value, nextRune)
					}
				case RUNE_QUOTE_DOUBLE:
					{
						state = STATE_APPEND
					}
				case RUNE_ESCAPE:
					{
						state = STATE_ESCAPING_QUOTED
					}
				default:
					{
						return nil, errors.New(fmt.Sprintf("Uknown rune: %v", nextRune))
					}
				}
			}
		case STATE_QUOTED: // in non-escaping single quotes
			{
				switch nextRuneType {
				case RUNE_EOF:
					{
						err = errors.New("EOF found when expecting closing quote.")
						break SCAN
					}
				case RUNE_CHAR, RUNE_SPACE, RUNE_QUOTE_DOUBLE, RUNE_ESCAPE,
					RUNE_PIPE, RUNE_REDIRECT:
					{
						value = append(value, nextRune)
					}
				case RUNE_QUOTE_SINGLE:
					{
						state = STATE_APPEND
					}
				default:
					{
						return nil, errors.New(fmt.Sprintf("Uknown rune: %v", nextRune))
					}
				}
			}
		case STATE_COMMENT:
			{
				switch nextRuneType {
				case RUNE_EOF:
					{
						break SCAN
					}
				case RUNE_CHAR, RUNE_QUOTE_DOUBLE, RUNE_ESCAPE,
					RUNE_QUOTE_SINGLE, RUNE_PIPE, RUNE_REDIRECT:
					{
						value = append(value, nextRune)
					}
				case RUNE_SPACE:
					{
						if nextRune == '\n' {
							state = STATE_START
							break SCAN
						} else {
							value = append(value, nextRune)
						}
					}
				default:
					{
						return nil, errors.New(fmt.Sprintf("Uknown rune: %v", nextRune))
					}
				}
			}
		default:
			{
				panic(fmt.Sprintf("Unexpected state: %v", state))
			}
		}
	}

	token := &Token{
		tokenType: tokenType,
		value:     TokenValue(value),
	}

	return token, err
}

// NextToken returns the next token in the stream, and an error value. If there are no more
// tokens available, the error value will be io.EOF.
func (t *Tokenizer) NextToken() (*Token, error) {
	return t.scanStream()
}
