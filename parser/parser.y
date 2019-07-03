%{
package parser

import (
  "github.com/plbalbi/json-example-generator/model"
  "log"
  "fmt"
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
  parsedStructField fieldAndDataType
}

%token <value> Identifier
%token TypeToken StructToken OpenCurlyBraceToken ClosingCurlyBraceToken ListTypeToken
%type <declaredStructsCount> StructDeclarations
%type <parsedStructField> Field
%type <newDataType> FieldType
%type <isList> ListOrNot

%start main

%%

main: StructDeclarations
{
    setResult(yylex, Result{
      structsCount: $1,
      typesRepository: dataTypeRepository,
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
{
  $$ = fieldAndDataType{
    name : $1,
    datatype : $2,
  }
}

FieldType: ListOrNot Identifier
{
  assembledDataTypeName := $2
  if $1 {
    assembledDataTypeName = fmt.Sprintf("[]%s", assembledDataTypeName)
  }
  if fromRepository := dataTypeRepository[assembledDataTypeName]; fromRepository != nil {
    $$ = fromRepository
    log.Printf("Just saw a datatype named: %s", fromRepository.GetName())
  } else {
    if $1 {
      if innerFromRepository := dataTypeRepository[$2]; innerFromRepository != nil {
        newListDataType := model.NewListDataType(assembledDataTypeName, innerFromRepository)
        dataTypeRepository[assembledDataTypeName] = newListDataType
        log.Printf("Just created new list datatype named %s", assembledDataTypeName)
        $$ = newListDataType
      } else {
        //Notify error
        log.Printf("Could not found inner datatype named: %s", $2)
      } 
    } else {
      //Notify error
      log.Printf("Not valid datatype named: %s", $2)
    }
  }
}

ListOrNot: {$$ = false}
  | ListTypeToken {$$ = true}
