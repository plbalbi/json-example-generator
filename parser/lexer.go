package parser

import (
	"log"
	"strings"
	"text/scanner"

	"github.com/golang-collections/collections/stack"
)

type lexer struct {
	result Result
	scan   scanner.Scanner
	states stack.Stack
}

type Result struct {
	structsCount int
}

type stateFunction func(*lexer) tokenType

type tokenType int

// Lex is somehow like the tokenStream.next() called it time it needs by the parser
func (lex *lexer) Lex(lval *yySymType) int {
	return int(lex.scanUntilTokenFound())
}

func (lex *lexer) Error(message string) {
}

func (lex *lexer) scanAndLog() {
	lex.scan.Scan()
	log.Printf("%s: %s", lex.scan.Position, lex.scan.TokenText())
}

func (lex *lexer) scanUntilTokenFound() tokenType {
	stateToRun := lex.states.Pop()
	if currentStateFunction, ok := stateToRun.(stateFunction); ok {
		return stateFunction(currentStateFunction)(lex)
	}
	return 0 // Horrible way, should return some other error instead of EOF
}

// Parse does the actual parsing
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

func structSignatureLookupState(lex *lexer) tokenType {
	lex.scanAndLog()
	switch tokenText := lex.scan.TokenText(); tokenText {
	case "type":
		// Push same state
		lex.states.Push(structSignatureLookupState)
		return TYPE_TOKEN
	case "struct":
		// Push same state
		lex.states.Push(structSignatureLookupState)
		return STRUCT_TOKEN
	default:
		return 0 // It seems '0' is recognized as an EOF token
	}
}
