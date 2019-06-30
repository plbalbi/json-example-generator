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
func Parse(input string) (Result, error) {
	return Result{
		structsCount: 0,
	}, nil
}
