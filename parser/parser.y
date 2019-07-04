%{
package parser

import (
  "github.com/plbalbi/json-example-generator/model"
  "log"
  "bytes"
)

func setResult(l yyLexer, v Result) {
  l.(*lexer).result = v
}

var GlobalRepository model.DataTypeRepository = model.GetDefaultDataTypeRepository()
var SeenDataTypes []string = make([]string, 0)
var Logger = log.Logger{}
var DeclaredStructs []string

func InitParser(){
  var logStream bytes.Buffer
	GlobalRepository = model.GetDefaultDataTypeRepository()
	SeenDataTypes = make([]string, 0)
  DeclaredStructs = make([]string, 0)
	Logger.SetOutput(&logStream)
}

%}

%union{
  value string
  newDataType string
  parsedStructField fieldAndDataType
  allParsedStructFields []fieldAndDataType
}

%token <value> Identifier
%token TypeToken StructToken OpenCurlyBraceToken ClosingCurlyBraceToken ListTypeToken
%type <value> StructOpening
%type <parsedStructField> Field
%type <newDataType> FieldType
%type <allParsedStructFields> StructFields

%start main

%%

main: StructDeclarations
{
    setResult(yylex, Result{
      declaredStructs: DeclaredStructs,
      typesRepository: GlobalRepository,
      })
}

StructDeclarations: StructDeclaration

StructDeclarations: StructDeclaration StructDeclarations

StructDeclaration: StructOpening StructFields ClosingCurlyBraceToken
{
  newStructName := $1
  newStruct := model.NewStructDataType(newStructName)
  for _,field := range $2 {
    // Already checked if struct fields datatype's exist
    newStruct.AddFieldNamed(field.name, field.datatypeName)
  }
  GlobalRepository[newStructName] = newStruct
  DeclaredStructs = append(DeclaredStructs, newStructName)
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

FieldType: Identifier
{
  $$ = $1
  Logger.Println("saw a simple type", $$)
  SeenDataTypes = append(SeenDataTypes, $$)
}
  | ListTypeToken FieldType
{
  $$ = "[]" + $2
  Logger.Println("saw a complex type", $$)
}