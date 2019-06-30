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

%token <value> Id
%token TypeDefinitionToken
%token StructDefinitionToken
%token StartStructDefinitionToken
%token EndStructDefinitionToken
%token IntTypeToken

%start Start

%%

Start: StructDefinitions

StructDefinitions:
    StructDefinition 
    | StructDefinition StructDefinitions

StructDefinition: 
    TypeDefinitionToken Id StructDefinitionToken StartStructDefinitionToken StructFields EndStructDefinitionToken

StructFields:
    StructField
    | StructField StructFields

StructField:
    Id Type

Type: IntTypeToken