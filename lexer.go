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
	"fmt"
	"io"
)

// Lexer turns an input stream in to a sequence of strings. Whitespace and
// comments are skipped.
type Lexer struct {
	tokenizer *Tokenizer
}

// Creates and returns a new lexer.
func NewLexer(reader io.Reader) (*Lexer, error) {
	tokenizer, err := NewTokenizer(reader)
	if err != nil {
		return nil, err
	}
	lexer := &Lexer{tokenizer: tokenizer}

	return lexer, nil
}

// NextWord return the next word, and an error value. If there are no more words, the error
// will be io.EOF.
func (lexer *Lexer) NextWord() (TokenValue, error) {
	var (
		token *Token
		err   error
	)

	for {
		token, err = lexer.tokenizer.NextToken()
		if err != nil {
			return "", err
		}
		switch token.tokenType {
		case TOKEN_WORD, TOKEN_PIPE, TOKEN_REDIRECT:
			{
				return token.value, nil
			}
		case TOKEN_COMMENT:
			{
				// skip comments
			}
		default:
			{
				panic(fmt.Sprintf("Unknown token type: %v", token.tokenType))
			}
		}
	}

	return "", io.EOF
}
