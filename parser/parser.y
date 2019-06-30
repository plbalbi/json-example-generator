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

%token <value> MemberIdentifier
%token <value> TypeIdentifier
%token <value> TypeOpening
%token TypeClosing
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

StructDeclaration: TypeOpening TypeClosing
