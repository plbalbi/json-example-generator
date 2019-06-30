%{
package parser

func setResult(l yyLexer, v Result) {
  l.(*lexer).result = v
}
%}

%union{
}

%token TYPE_TOKEN STRUCT_TOKEN START_STRUCT_DECL_TOKEN END_STRUCT_DECL_TOKEN
%token <value> ID

%start main

%%

main: StructDeclaration
{
    setResult(yylex, Result {
       structsCount: 1,
    })
}

StructDeclaration: TYPE_TOKEN ID STRUCT_TOKEN START_STRUCT_DECL_TOKEN END_STRUCT_DECL_TOKEN
