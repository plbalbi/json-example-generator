%{
package parser

import (
  "github.com/json-example-generator/model"
  "log"
)

func setResult(l yyLexer, v Result) {
  l.(*lexer).result = v
}

var dataTypeRepository model.DataTypeRepository = model.GetDefaultDataTypeRepository()

%}

%union{
  declaredStructsCount int
  value string
  newDataType model.DataType
  isList bool
}

%token <value> Identifier
%token TypeToken StructToken OpenCurlyBraceToken ClosingCurlyBraceToken ListTypeToken
%type <declaredStructsCount> StructDeclarations
%type <newDataType> FieldType
%type <isList> ListOrNot

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

StructDeclaration: StructOpening StructFields ClosingCurlyBraceToken

StructOpening: TypeToken Identifier StructToken OpenCurlyBraceToken

StructFields:  
  | Field StructFields

Field: Identifier FieldType

FieldType: ListOrNot Identifier
{
  if fromRepository := dataTypeRepository[$2]; fromRepository != nil {
    $$ = fromRepository
    log.Printf("Just saw a datatype named: %s", fromRepository.GetName())
  } else {
    log.Printf("Unrecognized datatype named: %s", $2)
  }
}

ListOrNot: {$$ = false}
  | ListTypeToken {$$ = true}
