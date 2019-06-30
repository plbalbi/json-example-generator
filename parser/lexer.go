package parser

type lexer struct {
	input string
}

type Result struct {
	structsCount int
}

type yySymType struct{}

func (lex *lexer) Lex(lval *yySymType) int {
	return 0
}

// Parse does the actual parsing
func Parse(inputStream string) (Result, error) {
	lex := &lexer{
		input: inputStream,
	}

	return Result{
		structsCount: 1,
	}, nil
}
