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

%token <value> Identifier
%token TypeToken StructToken OpenCurlyBraceToken ClosingCurlyBraceToken
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

StructDeclaration: TypeOpening TypeMembers ClosingCurlyBraceToken

TypeOpening: TypeToken Identifier StructToken OpenCurlyBraceToken

TypeMembers:  
  | TypeMember TypeMembers

TypeMember: Identifier Identifier
