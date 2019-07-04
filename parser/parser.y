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

var SeenDataTypes []string = make([]string, 0)
var globalRepository model.DataTypeRepository = model.GetDefaultDataTypeRepository()
var logger = log.Logger{}
var declaredStructs []string
var logStream bytes.Buffer
var freshIdentifier int
var structDependencyGraph map[string][]string

func InitParser(){
  logStream = bytes.Buffer{}
	globalRepository = model.GetDefaultDataTypeRepository()
	SeenDataTypes = make([]string, 0)
  declaredStructs = make([]string, 0)
  structDependencyGraph = make(map[string][]string, 0)
	logger.SetOutput(&logStream)
}

func RegisterNewStruct(name string, fields []fieldAndDataType){
  newStruct := model.NewStructDataType(name)
  for _, field := range fields {
    newStruct.AddFieldNamed(field.name, field.datatypeName)
  }
  globalRepository[name] = newStruct
}

func generateFreshIdentifier() string{
  for {
    freshIdentifier += 1
    id := "_f" + strconv.Itoa(freshIdentifier)
    if globalRepository[id] == nil{
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
      declaredStructs: declaredStructs,
      typesRepository: globalRepository,
      logRegistry: logStream.String(),
      structDependencyGraph: structDependencyGraph,
    })
}

StructDeclarations: StructDeclaration

StructDeclarations: StructDeclaration StructDeclarations

StructDeclaration: TypeName InlineStructDeclaration
{
  newStructName := $1
  RegisterNewStruct(newStructName, $2)
  declaredStructs = append(declaredStructs, newStructName)
  var adyacents []string
  for _, field := range $2 {
    adyacents = append(adyacents, field.datatypeName)
  }
  structDependencyGraph[newStructName] = adyacents
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
  logger.Println("saw a simple type", $$)
  SeenDataTypes = append(SeenDataTypes, $$)
}
  | ListTypeToken FieldType
{
  $$ = "[]" + $2
  logger.Println("saw a complex type", $$)
}
 | InlineStructDeclaration
{
  $$ = generateFreshIdentifier()
  RegisterNewStruct($$, $1)
}