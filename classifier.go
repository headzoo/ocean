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

// RuneType is the type of a UTF-8 character; a character, quote, space, escape.
type RuneType string

// RuneTypeMap is a map of RuneTokeType values.
type RuneTypeMap map[rune]RuneType

const (
	CLASS_CHAR              string = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789#._-,/@$*()+=:;&^%~"
	CLASS_SPACE             string = " \t\r\n"
	CLASS_ESCAPING_QUOTE    string = "\""
	CLASS_NONESCAPING_QUOTE string = "'"
	CLASS_ESCAPE            string = "\\"
	CLASS_PIPE              string = "|"
	CLASS_REDIRECT          string = "><"

	RUNE_UNKNOWN      RuneType = "UNKNOWN"
	RUNE_CHAR         RuneType = "CHAR"
	RUNE_SPACE        RuneType = "SPACE"
	RUNE_QUOTE_DOUBLE RuneType = "QUOTE_DOUBLE"
	RUNE_QUOTE_SINGLE RuneType = "QUOTE_SINGLE"
	RUNE_ESCAPE       RuneType = "ESCAPE"
	RUNE_PIPE         RuneType = "PIPE"
	RUNE_REDIRECT     RuneType = "REDIRECT"
	RUNE_EOF          RuneType = "EOF"
)

// Classifier classifies runes by type. This allows for different sorts of
// classifiers - those accepting extended non-ascii chars, or strict posix
// compatibility, for example.
type Classifier struct {
	typeMap RuneTypeMap
}

// Create and returns a new rune classifier.
func NewClassifier() *Classifier {
	classifier := &Classifier{
		typeMap: make(RuneTypeMap),
	}

	classifier.AddClassification(CLASS_CHAR, RUNE_CHAR)
	classifier.AddClassification(CLASS_SPACE, RUNE_SPACE)
	classifier.AddClassification(CLASS_ESCAPING_QUOTE, RUNE_QUOTE_DOUBLE)
	classifier.AddClassification(CLASS_NONESCAPING_QUOTE, RUNE_QUOTE_SINGLE)
	classifier.AddClassification(CLASS_ESCAPE, RUNE_ESCAPE)
	classifier.AddClassification(CLASS_PIPE, RUNE_PIPE)
	classifier.AddClassification(CLASS_REDIRECT, RUNE_REDIRECT)

	return classifier
}

// Classify returns the rune token type.
func (classifier *Classifier) Classify(r rune) RuneType {
	return classifier.typeMap[r]
}

// addRuneClass registers a rune and it's classification.
func (classifier *Classifier) AddClassification(runes string, runeType RuneType) {
	for _, r := range runes {
		classifier.typeMap[r] = runeType
	}
}
