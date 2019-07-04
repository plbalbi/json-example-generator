package parser

import (
	"errors"
	"log"
	"testing"

	"github.com/plbalbi/json-example-generator/model"
	"github.com/stretchr/testify/assert"
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
			"missing type declarations fails",
			"type perro struct { hola perro que haces como va }",
			errors.New("Type 'haces' was not declared"),
			func(result *Result) bool { return model.CountStructDataTypes(result.typesRepository) == 1 },
		},
		{
			"single structs single list field is parsed correctly",
			"type perro struct { hola []int }",
			nil,
			func(result *Result) bool { return model.CountStructDataTypes(result.typesRepository) == 1 },
		},
		{
			"slices are not registered in repository ",
			"type gato struct { hola []string }",
			errors.New("asd"),
			func(result *Result) bool { return result.typesRepository["[]string"] == nil },
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
			if testCase.expectedError != nil {
				assert.EqualError(t, err, testCase.expectedError.Error())
			} else {
				t.Errorf("Expected to parse input correctly, but got this error: %s", err.Error())
			}
		}
		//fmt.Println(result.typesRepository)
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

//TODO: Should find a way to display struct field in the order they are defined.
func TestRandomJsonGeneration(t *testing.T) {
	result, err := Parse(`type test struct {
		nombre string
		edad int
		gustosDeHelado []string
	}`)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	//fmt.Println(result.typesRepository["test"].Generate(result.typesRepository))
	//t.Logf("Parser got:\n%s", result.typesRepository["test"].Generate())
}

func Test00(t *testing.T) {
	result, err := Parse(`
	type test struct {
		nombre string
		edad int
		gustosDeHelado []gusto
	}
	`)
	assert.Error(t, err)
	assert.NotNil(t, result)
	//t.Logf("Parser got:\n%s", result.typesRepository["test"].Generate())
}

func Test02(t *testing.T) {
	result, err := Parse(`
	type test struct {
		nombre string
		edad int
		gustosDeHelado []gusto
	}
	type gusto struct {
		nombre string
		granizado bool
	}
	`)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	//fmt.Println(result.typesRepository)
	//fmt.Println(result.typesRepository["test"].Generate(result.typesRepository))
	//t.Logf("Parser got:\n%s", result.typesRepository["test"].Generate())
}

func Test03(t *testing.T) {
	result, err := Parse(`
	type test struct {
		nombre string
		edad int
		gustosDeHelado []gusto
	}
	type gusto struct {
		nombre string
		granizado bool
	}
	`)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	//fmt.Println(result.typesRepository)
	//fmt.Println(result.typesRepository["test"].Generate(result.typesRepository))
	//t.Logf("Parser got:\n%s", result.typesRepository["test"].Generate())
}
