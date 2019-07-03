package parser

import (
	"log"
	"testing"

	"github.com/plbalbi/json-example-generator/model"
)

func TestParser(t *testing.T) {
	tests := []simpleParserTestCase{
		{
			"single struct is parser correctly",
			"type perro struct { }",
			nil,
			func(result *Result) bool { return model.CountStructDataTypes(result.typesRepository) == 1 },
		},
		{
			"two structs are parsed correctly",
			"type perro struct { } type gato struct { }",
			nil,
			func(result *Result) bool {
				return model.CountStructDataTypes(result.typesRepository) == 2 &&
					result.typesRepository["perro"] != nil &&
					result.typesRepository["gato"] != nil
			},
		},
		{
			"single structs with one field is parsed correctly",
			"type perro struct { hola perro }",
			nil,
			func(result *Result) bool { return model.CountStructDataTypes(result.typesRepository) == 1 },
		},
		{
			"single structs with three field is parsed correctly",
			"type perro struct { hola perro que haces como va }",
			nil,
			func(result *Result) bool { return model.CountStructDataTypes(result.typesRepository) == 1 },
		},
		{
			"single structs single list field is parsed correctly",
			"type perro struct { hola []int }",
			nil,
			func(result *Result) bool { return model.CountStructDataTypes(result.typesRepository) == 1 },
		},
		{
			"struct with list type gets it registered in repository",
			"type gato struct { hola []string }",
			nil,
			func(result *Result) bool { return result.typesRepository["[]string"] != nil },
		},
		{
			"struct with non existing inner type omits registering corresponding list type",
			"type gato struct { hola []falopa }",
			nil,
			func(result *Result) bool { return result.typesRepository["[]falopa"] == nil },
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.description, generateSingleParserTest(testCase))
	}
}

type simpleParserTestCase struct {
	description     string
	input           string
	expectedError   error
	resultPredicate func(*Result) bool
}

func generateSingleParserTest(testCase simpleParserTestCase) func(*testing.T) {
	return func(t *testing.T) {
		result, err := Parse(testCase.input)
		if err != nil {
			if testCase.expectedError != nil && err.Error() != testCase.expectedError.Error() {
				t.Errorf("An error with message '%s' was expected, but got '%s'", testCase.expectedError.Error(), err.Error())
			} else {
				t.Errorf("Expected to parse input correctly, but got this error: %s", err.Error())
			}
		}
		if !testCase.resultPredicate(&result) {
			t.Errorf("Failed to evaluate test predicate")
		}
	}
}

func TestMapStuff(t *testing.T) {
	repository := model.GetDefaultDataTypeRepository()
	nonExistingDataType := repository["caca"]
	if nonExistingDataType != nil {
		t.Errorf("Expected repository to return nil on non existing data type")
		log.Print("hola")
	}
}
