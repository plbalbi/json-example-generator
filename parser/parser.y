%{
package parser

import (
  "github.com/plbalbi/json-example-generator/model"
  "log"
  "bytes"
  "strconv"
)

func setResult(l yyLexer, v Result) {
  l.(*lexer).result = v
}

var GlobalRepository model.DataTypeRepository = model.GetDefaultDataTypeRepository()
var SeenDataTypes []string = make([]string, 0)
var Logger = log.Logger{}
var DeclaredStructs []string
var LogStream bytes.Buffer
var freshIdentifier int

func InitParser(){
  LogStream = bytes.Buffer{}
	GlobalRepository = model.GetDefaultDataTypeRepository()
	SeenDataTypes = make([]string, 0)
  DeclaredStructs = make([]string, 0)
	Logger.SetOutput(&LogStream)
}

func RegisterNewStruct(name string, fields []fieldAndDataType){
  newStruct := model.NewStructDataType(name)
  for _, field := range fields {
    newStruct.AddFieldNamed(field.name, field.datatypeName)
  }
  GlobalRepository[name] = newStruct
}

func generateFreshIdentifier() string{
  for {
    freshIdentifier += 1
    id := "_f" + strconv.Itoa(freshIdentifier)
    if GlobalRepository[id] == nil{
      return id
    }
  }
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
%type <parsedStructField> Field
%type <newDataType> FieldType
%type <allParsedStructFields> StructFields
%type <allParsedStructFields> InlineStructDeclaration
%type <newDataType> TypeName

%start main

%%

main: StructDeclarations
{
    setResult(yylex, Result{
      declaredStructs: DeclaredStructs,
      typesRepository: GlobalRepository,
      logRegistry: LogStream.String(),
    })
}

StructDeclarations: StructDeclaration

StructDeclarations: StructDeclaration StructDeclarations

StructDeclaration: TypeName InlineStructDeclaration
{
  newStructName := $1
  RegisterNewStruct(newStructName, $2)
  DeclaredStructs = append(DeclaredStructs, newStructName)
}

TypeName: TypeToken Identifier
{
  $$ = $2
}

InlineStructDeclaration: StructToken OpenCurlyBraceToken StructFields ClosingCurlyBraceToken
{
  $$ = $3
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
 | InlineStructDeclaration
{
  $$ = generateFreshIdentifier()
  RegisterNewStruct($$, $1)
}