%{
package parser

func setResult(l yyLexer, v Result) {
  l.(*lexer).result = v
}
%}

%union{
  result Result
  part string
  ch byte
}

%token <value> ID

%start Start

%%

Start: StructDefinitions

StructDefinitions:
    StructDefinition 
    | StructDefinition StructDefinitions

StructDefinition: 
    'type' ID 'struct {' StructFields '}'

StructFields:
    StructField
    | StructField StructFields

StructField:
    ID Type

Type: 'int'