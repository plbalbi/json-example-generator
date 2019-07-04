package parser

import (
	"errors"
	"log"
	"regexp"
	"strings"
	"text/scanner"

	"github.com/golang-collections/collections/queue"
	"github.com/plbalbi/json-example-generator/model"
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
	declaredStructs       []string
	typesRepository       model.DataTypeRepository
	logRegistry           string
	structDependencyGraph map[string][]string
}

func (res *Result) StructsCount() int {
	// return res.structsCount
	return model.CountStructDataTypes(res.typesRepository)
}

func (res *Result) GetDataTypeNames() []string {
	keys := make([]string, len(res.typesRepository))

	i := 0
	for k := range res.typesRepository {
		keys[i] = k
		i++
	}
	return keys
}

func (res *Result) FirstDataTypeSeen() string {
	return res.declaredStructs[0]
}

func (res *Result) GenerateDataType() string {
	return res.typesRepository[res.FirstDataTypeSeen()].Generate(res.typesRepository)
}

//Parse lexes and parses the file and returns the parsed text.
func Parse(inputStream string) (Result, error) {
	lex := newLexer(inputStream)
	//Clearing global repository between parse calls
	InitParser()
	yyParse(lex)

	// Check if all seen data types were defined
	for _, typeName := range SeenDataTypes {
		if lex.result.typesRepository[typeName] == nil {
			return lex.result, errors.New("Type '" + typeName + "' was not declared")
		}
	}

	// Check if a type definition was seen more than once
	seenTypeDeclarations := make(map[string]bool)
	for _, typeName := range lex.result.declaredStructs {
		if seenTypeDeclarations[typeName] {
			return lex.result, errors.New("Multiple declarations of type '" + typeName + "'")
		}
		seenTypeDeclarations[typeName] = true
	}

	// Look for circular definitons
	for typeName, _ := range lex.result.structDependencyGraph {
		if reachesSelf(typeName, structDependencyGraph) {
			return lex.result, errors.New("Circular definition of type '" + typeName + "'")
		}
	}
	return lex.result, lex.err
}

func reachesSelf(typeName string, structDependencyGraph map[string][]string) bool {
	if structDependencyGraph[typeName] == nil {
		return false
	}
	// DFS
	seen := make(map[string]bool)
	var lefts queue.Queue
	lefts.Enqueue(typeName)
	for lefts.Len() > 0 {
		current := lefts.Dequeue().(string)
		if seen[current] {
			return true
		} else {
			seen[current] = true
			for _, neighbour := range structDependencyGraph[current] {
				// We only care about structs, it may not be one
				if structDependencyGraph[neighbour] != nil {
					lefts.Enqueue(neighbour)
				}
			}
		}
	}
	return false
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
	identifierMatcher := regexp.MustCompile(`^[a-zA-Z,0-9,_]+$`).MatchString
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
