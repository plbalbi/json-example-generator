%{
package parser

func setResult(l yyLexer, v Result) {
  l.(*lexer).result = v
}
%}

%union{
  declaredStructsCount int
}

%token TYPE_TOKEN STRUCT_TOKEN START_STRUCT_DECL_TOKEN END_STRUCT_DECL_TOKEN
%token <value> ID
%type <declaredStructsCount> StructDeclarations

%start main

%%

main: StructDeclarations
{
    setResult(yylex, Result{
      structsCount: $1,
      })
}

StructDeclarations: StructDeclaration
{
  $$ = 1
}

StructDeclarations: StructDeclaration StructDeclarations
{
  $$ = $2 + 1
}

StructDeclaration: TYPE_TOKEN
