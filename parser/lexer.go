package parser

import (
	"strings"
	"unicode/utf8"
)

const eof = -1

type lexer struct {
	input       string
	start       int
	pos         int
	result      Result
	tokenStream chan item
	state       stateFunction
	width       int
}

type itemType int

type item struct {
	typ   itemType
	value string
}

type Result struct {
	structsCount int
}

type stateFunction func(*lexer) stateFunction

// Lex is somehow like the tokenStream.next() called it time it needs by the parser
func (lex *lexer) Lex(lval *yySymType) int {
	for {
		select {
		case someItem := <-lex.tokenStream:
			return int(someItem.typ)
		default:
			lex.state = lex.state(lex)
		}
	}
	return 0
}
func (lex *lexer) emitToken(tokenType itemType) {
	lex.tokenStream <- item{
		typ:   tokenType,
		value: lex.input[lex.start:lex.pos],
	}
	lex.start = lex.pos
}

func (lex *lexer) Error(message string) {
}

// Parse does the actual parsing
func Parse(inputStream string) (Result, error) {
	lex := newLexer(inputStream)
	yyParse(lex)
	return lex.result, nil
}

func newLexer(inputStream string) *lexer {
	return &lexer{
		input:       inputStream,
		start:       0,
		pos:         0,
		tokenStream: make(chan item, 2), // Item channel has buffering
		state:       lookingForStructDefinitionState,
	}
}

// Returns a whitespace consuming state, which when a non-whitespace character is reached, will continue into nextState
func whiteSpaceConsumerWithNext(nextState stateFunction) stateFunction {
	return func(lex *lexer) stateFunction {
		for {
			seenCharacter := lex.next()
			if seenCharacter == -1 {
				return nil
			} else if seenCharacter != ' ' {
				return nextState
			}
		}
	}
}

func (lex *lexer) next() rune {
	if lex.pos > len(lex.input) {
		lex.width = 0
		return eof
	}
	var r rune
	r, lex.width = utf8.DecodeRuneInString(lex.input[lex.pos:])
	lex.pos += lex.width
	return r

}

const typeItem = "type"
const structItem = "struct"

func lookingForStructDefinitionState(lex *lexer) stateFunction {
	for {
		if strings.HasPrefix(lex.input[lex.pos:], typeItem) {
			lex.start = lex.pos
			return typeLiteralParsingState
		}
		if lex.next() == eof {
			break
		}
	}
	// Maybe emit an EOF token
	return nil
}

func typeLiteralParsingState(lex *lexer) stateFunction {
	lex.pos += len(typeItem)
	lex.emitToken(TYPE_TOKEN)
	return nil
}
