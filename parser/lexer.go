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

// Lex is somehow like the tokenStream.next() called it time it needs by the parser
func (lex *lexer) Lex(lval *yySymType) int {
	lex.scan.Scan()
	log.Printf("%s: %s", lex.scan.Position, lex.scan.TokenText())
	switch tokenText := lex.scan.TokenText(); tokenText {
	case "type":
		return TYPE_TOKEN
	case "struct":
		return STRUCT_TOKEN
	default:
		return 0 // It seems '0' is recognized as an EOF token
	}
}

func (lex *lexer) Error(message string) {
}

// Parse does the actual parsing
func Parse(inputStream string) (Result, error) {
	lex := newLexer(inputStream)
	lex.scan.Init(strings.NewReader(inputStream))
	yyParse(lex)
	return lex.result, nil
}

func newLexer(inputStream string) *lexer {
	return &lexer{}
}
