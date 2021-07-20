package main

import (
	"strings"
	"text/scanner"
	"unicode"
)

var replacementLetters = "abcdefghijklmnopqrstuvxyz"

type shortNames struct {
	set map[string]bool
}

func NewShortNames() *shortNames {
	return &shortNames{make(map[string]bool)}
}

// TODO use multiple characters if all single characters are taken
func (shortNames *shortNames) nextShortName(providedWords ...string) (string, bool) {
	// add defaults as additonal word
	words := append(providedWords, replacementLetters)

	// iterate over all characters in all provided words
	for _, word := range words {
		for _, char := range word {
			newName := strings.ToLower(string(char))
			_, ok := shortNames.set[newName]
			// find a character that is not yet in use
			if !ok {
				// add new value to set
				shortNames.set[newName] = true
				return newName, true
			}
		}
	}

	return "", false
}

func isVariable(token string) bool {
	return token == "var"
}

func isSpaceString(token string) bool {
	return token == " "
}

func isWhiteSpace(token string) bool {
	if len(token) == 1 {
		runes := []rune(token)
		return unicode.IsSpace(runes[0])
	}
	return false
}

func isWord(whiteSpaceOrWord string) bool {
	isSpace := isWhiteSpace(whiteSpaceOrWord)
	isEmpty := len(whiteSpaceOrWord) == 0

	return !isSpace && !isEmpty
}
func handleShortAssignment(renamedVariables map[string]string, shortNames *shortNames, buffer *buffer) {
	// bla :=
	// buffer: [... "bla" , ":", "="]
	// or:
	// buffer: [... "bla", " ", ":", "="]
	colon, _ := buffer.Read(1)

	// the = could also be an assignment only or part of a comparison
	// we are only interested if there is a : before
	if colon == ":" {
		// prev character can be either whitespace or the variable name
		whiteSpaceOrToken, _ := buffer.Read(2)
		offset := 2
		if !isWord(whiteSpaceOrToken) {
			whiteSpaceOrToken, _ = buffer.Read(3)
			offset = 3
		}
		if isWord(whiteSpaceOrToken) {
			newShortName, shortNameFound := shortNames.nextShortName(whiteSpaceOrToken)
			if shortNameFound {
				renamedVariables[whiteSpaceOrToken] = newShortName
				buffer.Replace(offset, newShortName)
			}
		}
	}
}

func handleVar(s *scanner.Scanner, usedVariableName map[string]string, shortNames *shortNames, buffer *buffer) string {
	s.Scan() // space
	s.Scan() // var name
	varName := s.TokenText()
	var overflow1, overflow2 string

	// only write a new entry for the variable if it is not already in the map
	_, variableExists := usedVariableName[varName]
	if !variableExists {
		// only write if we found a new shortName
		newShortName, shortNameFound := shortNames.nextShortName(varName)
		if shortNameFound {
			usedVariableName[varName] = newShortName
			overflow1, _ = buffer.Skim()
			buffer.Push(" ")
			overflow2, _ = buffer.Skim()
			buffer.Push(newShortName)
		}
	}
	return overflow1 + overflow2
}

func NewScanner(input string) *scanner.Scanner {
	var s scanner.Scanner
	s.Init(strings.NewReader(input))
	// scanner ignores comments, so those will be removed!
	// don't skip any whitespace while scanning
	s.Whitespace ^= 1<<'\t' | 1<<'\n' | 1<<'\v' | 1<<'\f' | 1<<'\r' | 1<<' '
	return &s
}

func minimize(code string) (minimizedCode string) {
	// Remove whitespace
	code = strings.Replace(code, "\n\n", "\n", 1)
	var renamedVariables = make(map[string]string)
	var shortNames = NewShortNames()

	s := NewScanner(code)

	// create the read buffer
	buffer := NewBuffer(4)

	for token := s.Scan(); token != scanner.EOF; token = s.Scan() {
		tokenText := s.TokenText()
		// check for multipe space tokens in a row and only keep one
		// removing multiple whitespace tokens is harder, because it may remove linebreaks with tabs after them :(
		prevToken, hasPrev := buffer.Read(0)
		if hasPrev && isSpaceString(tokenText) && isSpaceString(prevToken) {
			buffer.Replace(0, tokenText)
		} else {
			buffer.Push(tokenText)
		}

		// is tokenText a word that needs to be replaced?
		replacement, exists := renamedVariables[tokenText]
		if exists {
			buffer.Replace(0, replacement)
		}

		// is it a variable declaration with var?
		if isVariable(tokenText) {
			overflow := handleVar(s, renamedVariables, shortNames, buffer)
			minimizedCode += overflow
		}
		// find variables assigned with :=
		if tokenText == "=" {
			handleShortAssignment(renamedVariables, shortNames, buffer)
		}

		// TODO find function params ?

		// TODO find function name?

		// TODO named function return values ?

		// if the buffer is full, we skim the first entry and write it
		tokenToWrite, bufferFull := buffer.Skim()
		if bufferFull {
			minimizedCode += tokenToWrite
		}
		// DEBUG
		// fmt.Printf("%s: %s\n", s.Position, tokenText)
	}
	// empty the buffer after everything was read
	for token, ok := buffer.Shift(); ok; token, ok = buffer.Shift() {
		minimizedCode += token
	}
	return
}
