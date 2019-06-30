package parser

import (
	"log"
	"strings"
	"text/scanner"

	"github.com/golang-collections/collections/stack"
)

type lexer struct {
	result         Result
	scan           scanner.Scanner
	states         stack.Stack
	currentSymType *yySymType
}

//Result is the object in which the parser transmits the parsed text.
type Result struct {
	structsCount int
}

//stateFunction is the base type for lexer states. Contains all the logic for lexing the text,
//communicating found values with the parser, and choosing the following states.
type stateFunction func(*lexer) tokenType

type tokenType int

//Lex is somehow like the tokenStream.next() called it time it needs by the parser
func (lex *lexer) Lex(currentSymType *yySymType) int {
	lex.currentSymType = currentSymType
	return int(lex.scanUntilTokenFound())
}

// TODO: Implement lexer/parser error handling.
func (lex *lexer) Error(message string) {
}

func (lex *lexer) scanAndLog() {
	lex.scan.Scan()
	log.Printf("%s: %s", lex.scan.Position, lex.scan.TokenText())
}

func (lex *lexer) scanUntilTokenFound() tokenType {
	stateToRun := lex.states.Pop()
	// Had to hardcode the stateFunction type in here, no way out
	if currentStateFunction, ok := stateToRun.(func(*lexer) tokenType); ok {
		// casting and calling current state
		return stateFunction(currentStateFunction)(lex)
	}
	return 0 // Horrible way, should return some other error instead of EOF
}

//Parse does the actual parsing
func Parse(inputStream string) (Result, error) {
	lex := newLexer(inputStream)
	yyParse(lex)
	return lex.result, nil
}

func newLexer(inputStream string) *lexer {
	brandNewLexer := &lexer{}
	brandNewLexer.states.Push(structSignatureLookupState)
	brandNewLexer.scan.Init(strings.NewReader(inputStream))
	return brandNewLexer
}

//structSignatureLookupState scans until the next type token is found.
//Then delegates to the ID scanner lexer state, identifierLexer.
func structSignatureLookupState(lex *lexer) tokenType {
	lex.scanAndLog()
	switch tokenText := lex.scan.TokenText(); tokenText {
	case "type":
		// Push same state
		lex.states.Push(identifierLexer)
		return TYPE_TOKEN
	default:
		return 0 // It seems '0' is recognized as an EOF token
	}
}

//identifierLexer scans the next token, which being after a 'type' symbol, is the id
// of the struct being defined. Also, transmits it to the parser through te SymType.
func identifierLexer(lex *lexer) tokenType {
	lex.scanAndLog()
	tokenText := lex.scan.TokenText()
	lex.currentSymType.value = tokenText
	lex.states.Push(structSignatureLookupState)
	return ID
}
