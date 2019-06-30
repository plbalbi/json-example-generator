package parser

import (
	"errors"
	"log"
	"regexp"
	"strings"
	"text/scanner"

	"github.com/golang-collections/collections/stack"
)

type lexer struct {
	result         Result
	err            error
	scan           scanner.Scanner
	states         stack.Stack
	currentSymType *yySymType
}

//Result is the object in which the parser transmits the parsed text.
type Result struct {
	structsCount int
}

// TODO: This could be concurrent. The lexer runs on one routine, feeding the parsed tokens into a channel,
// while the parser consumes from there. Some buffering can be added, so the lexer does not overfill. Also,
// to communicate the values attached to certain tokens, some kind of tokenWithAttachment struct could be used,
// and the Lex methods consumes from the channel, and separates the tokenType from the tokenAttachment.

//Lex is somehow like the tokenStream.next() called it time it needs by the parser.
func (lex *lexer) Lex(currentSymType *yySymType) int {
	lex.currentSymType = currentSymType
	lex.scan.Scan()
	switch scannedTokenText := lex.scan.TokenText(); scannedTokenText {
	case "type":
		return TypeToken
	case "struct":
		return StructToken
	case "{":
		return OpenCurlyBraceToken
	case "}":
		return ClosingCurlyBraceToken
	default:
		// TODO: Change this regex for sth else. Too overkill!
		identifierMatcher := regexp.MustCompile(`^[a-zA-Z]+$`).MatchString
		if !identifierMatcher(scannedTokenText) {
			return 0
		}
		lex.currentSymType.value = scannedTokenText
		return Identifier
	}
}

// TODO: Implement lexer/parser error handling.
func (lex *lexer) Error(message string) {
	lex.err = errors.New(message)
}

// TODO: Log only on some debug mode?
func (lex *lexer) scanAndLog() string {
	lex.scan.Scan()
	log.Printf("%s: %s", lex.scan.Position, lex.scan.TokenText())
	return lex.scan.TokenText()
}

//Parse lexes and parses the file and returns the parsed text.
func Parse(inputStream string) (Result, error) {
	lex := newLexer(inputStream)
	yyParse(lex)
	return lex.result, lex.err
}

//newLexer creates a brand new lexer object.
func newLexer(inputStream string) *lexer {
	brandNewLexer := &lexer{}
	brandNewLexer.scan.Init(strings.NewReader(inputStream))
	return brandNewLexer
}
