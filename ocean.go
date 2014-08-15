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
	"errors"
	"fmt"
)

// Split splits a string in to a slice of strings, based upon shell-style rules for
// quoting, escaping, and spaces.
func Split(s string) ([]TokenValue, error) {
	lexer := NewLexer(strings.NewReader(s))
	subStrings := []TokenValue{}

	for {
		word, err := lexer.NextWord()
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

// errorf returns a new error with a formatted message.
func errorf(format string, a... interface{}) error {
	return errors.New(fmt.Sprintf(format, a...))
}
