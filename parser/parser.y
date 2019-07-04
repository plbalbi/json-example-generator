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

var GlobalRepository model.DataTypeRepository = model.GetDefaultDataTypeRepository()

%}

%union{
  declaredStructsCount int
  value string
  newDataType string
  isList bool
  parsedStructField fieldAndDataType
  allParsedStructFields []fieldAndDataType
}

%token <value> Identifier
%token TypeToken StructToken OpenCurlyBraceToken ClosingCurlyBraceToken ListTypeToken
%type <value> StructOpening
%type <declaredStructsCount> StructDeclarations
%type <parsedStructField> Field
%type <newDataType> FieldType
%type <isList> ListOrNot
%type <allParsedStructFields> StructFields

%start main

%%

main: StructDeclarations
{
    setResult(yylex, Result{
      structsCount: $1,
      typesRepository: GlobalRepository,
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
{
  newStructName := $1
  newStruct := model.NewStructDataType(newStructName)
  for _,field := range $2 {
    // Already checked if struct fields datatype's exist
    newStruct.AddFieldNamed(field.name, field.datatypeName)
  }
  GlobalRepository[newStructName] = newStruct
}

StructOpening: TypeToken Identifier StructToken OpenCurlyBraceToken
{
  //Pass through struct name
  //TODO: Maybe fail here if struct already defined?
  $$ = $2
}

StructFields:  { $$ = make([]fieldAndDataType, 0) }
  | Field StructFields { $$ = append($2, $1) }

Field: Identifier FieldType
{
  $$ = fieldAndDataType{
    name : $1,
    datatypeName : $2,
  }
}

FieldType: ListOrNot Identifier
{
  assembledDataTypeName := $2
  if $1 {
    assembledDataTypeName = fmt.Sprintf("[]%s", assembledDataTypeName)
  }
  $$ = assembledDataTypeName
  if fromRepository := GlobalRepository[assembledDataTypeName]; fromRepository != nil {
    //$$ = fromRepository
    log.Printf("Just saw a datatype named: %s", fromRepository.GetName())
  } else {
    if $1 {
      if innerFromRepository := GlobalRepository[$2]; innerFromRepository != nil {
        newListDataType := model.NewListDataType(assembledDataTypeName, innerFromRepository)
        GlobalRepository[assembledDataTypeName] = newListDataType
        log.Printf("Just created new list datatype named %s", assembledDataTypeName)
        //$$ = newListDataType
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
