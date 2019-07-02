package parser

import (
	"errors"
	"log"
	"regexp"
	"strings"
	"text/scanner"

	"github.com/golang-collections/collections/queue"
)

type lexer struct {
	result         Result
	err            error
	scan           scanner.Scanner
	currentSymType *yySymType
	itemLeftToEmit queue.Queue
}

const errorToken = 0

type lexedItem struct {
	itemType int
	value    string
}

//Result is the object in which the parser transmits the parsed text.
type Result struct {
	structsCount int
}

//Parse lexes and parses the file and returns the parsed text.
func Parse(inputStream string) (Result, error) {
	lex := newLexer(inputStream)
	yyParse(lex)
	return lex.result, lex.err
}

// TODO: This could be concurrent. The lexer runs on one routine, feeding the parsed tokens into a channel,
// while the parser consumes from there. Some buffering can be added, so the lexer does not overfill. Also,
// to communicate the values attached to certain tokens, some kind of tokenWithAttachment struct could be used,
// and the Lex methods consumes from the channel, and separates the tokenType from the tokenAttachment.

//Lex is somehow like the tokenStream.next() called it time it needs by the parser.
func (lex *lexer) Lex(currentSymType *yySymType) int {
	if lex.itemLeftToEmit.Len() == 0 {
		lex.doLex()
	}
	// Pop oldest enqueued item
	oldestLexedItem := lex.itemLeftToEmit.Dequeue().(lexedItem)
	currentSymType.value = oldestLexedItem.value
	return oldestLexedItem.itemType
}

// TODO: Implement lexer/parser error handling.
func (lex *lexer) Error(message string) {
	lex.err = errors.New(message)
}

func (lex *lexer) doLex() {
	lex.scan.Scan()
	switch scannedTokenText := lex.scan.TokenText(); scannedTokenText {
	case "type":
		lex.emitItemOfType(TypeToken)
	case "struct":
		lex.emitItemOfType(StructToken)
	case "{":
		lex.emitItemOfType(OpenCurlyBraceToken)
	case "}":
		lex.emitItemOfType(ClosingCurlyBraceToken)
	case "[":
		// List type definition
		lex.scan.Scan()
		if lex.scan.TokenText() != "]" {
			lex.emitItemOfType(errorToken)
		}
		lex.emitItemOfType(ListTypeToken)
		lex.scan.Scan()
		lex.lexIdentifier()
	default:
		lex.lexIdentifier()
	}
}

func (lex *lexer) lexIdentifier() {
	scannedTokenText := lex.scan.TokenText()
	// TODO: Change this regex for sth else. Too overkill!
	if strings.HasPrefix(scannedTokenText, "[]") {
		// It's a list type
		lex.emitItemOfType(ListTypeToken)
		scannedTokenText = scannedTokenText[2:]
	}
	identifierMatcher := regexp.MustCompile(`^[a-zA-Z]+$`).MatchString
	if !identifierMatcher(scannedTokenText) {
		lex.emitItemOfType(errorToken)
	}
	lex.emitItemOfType(Identifier)
}

func (lex *lexer) emitItemOfType(emittedItemType int) {
	lex.itemLeftToEmit.Enqueue(lexedItem{
		itemType: emittedItemType,
		value:    lex.scan.TokenText(),
	})
}

// TODO: Log only on some debug mode?
func (lex *lexer) scanAndLog() string {
	lex.scan.Scan()
	log.Printf("%s: %s", lex.scan.Position, lex.scan.TokenText())
	return lex.scan.TokenText()
}

//newLexer creates a brand new lexer object.
func newLexer(inputStream string) *lexer {
	brandNewLexer := &lexer{}
	brandNewLexer.scan.Init(strings.NewReader(inputStream))
	return brandNewLexer
}
