%{
package parser

func setResult(l yyLexer, v Result) {
  l.(*lexer).result = v
}
%}

%union{
  declaredStructsCount int
  value string
}

%token TYPE_TOKEN STRUCT_TOKEN START_STRUCT_DECL_TOKEN END_STRUCT_DECL_TOKEN
%token <value> ID
%token <value> TypeOpening
%type <declaredStructsCount> StructDeclarations

%start main

%%

main: StructDeclarations
{
    setResult(yylex, Result{
      structsCount: $1,
      })
}

StructDeclarations: TypeOpening
{
  $$ = 1
}

StructDeclarations: StructDeclaration StructDeclarations
{
  $$ = $2 + 1
}

StructDeclaration: TypeOpening
