package parser

import (
	"errors"
	"log"
	"testing"

	"github.com/plbalbi/json-example-generator/model"
	"github.com/stretchr/testify/assert"
)

func TestParser(t *testing.T) {
	tests := []struct {
		testDescription            string
		input                      string
		expectedError              error
		expectedStructDeclarations []string
		resultPredicate            func(*Result) bool
	}{
		{
			"single struct is parser correctly",
			"type perro struct { }",
			nil,
			[]string{"perro"},
			func(result *Result) bool { return model.CountStructDataTypes(result.typesRepository) == 1 },
		},
		{
			"two structs are parsed correctly",
			"type perro struct { } type gato struct { }",
			nil,
			[]string{"perro", "gato"},
			func(result *Result) bool {
				return model.CountStructDataTypes(result.typesRepository) == 2 &&
					result.typesRepository["perro"] != nil &&
					result.typesRepository["gato"] != nil
			},
		},
		{
			// TODO: should this fail?
			"single structs with one field is parsed correctly",
			"type perro struct { hola perro }",
			nil,
			[]string{"perro"},
			func(result *Result) bool { return model.CountStructDataTypes(result.typesRepository) == 1 },
		},
		{
			"single structs single list field is parsed correctly",
			"type perro struct { hola []int }",
			nil,
			[]string{"perro"},
			func(result *Result) bool { return model.CountStructDataTypes(result.typesRepository) == 1 },
		},
		{
			"slices are not registered in repository",
			"type gato struct { hola []string }",
			nil,
			[]string{"gato"},
			func(result *Result) bool { return result.typesRepository["[]string"] == nil },
		},
		{
			"missing type declarations fails",
			`type test struct {
				nombre string
				edad int
				gustosDeHelado []gusto
			}`,
			errors.New("Type 'gusto' was not declared"),
			[]string{"test"},
			func(result *Result) bool { return result.typesRepository["[]string"] == nil },
		},
		{
			"declaring type later it valid",
			`type test struct {
				nombre string
				edad int
				gustosDeHelado []gusto
			}
			type gusto struct {
				nombre string
				granizado bool
			}`,
			nil,
			[]string{"test", "gusto"},
			func(result *Result) bool { return result.typesRepository["[]string"] == nil },
		},
		{
			"slices can be nested",
			`type test struct {
				nombre string
				edad int
				gustosDeHelado [][][][][]gusto
			}
			type gusto struct {
				nombre []string
				granizado bool
			}`,
			nil,
			[]string{"test", "gusto"},
			func(result *Result) bool { return result.typesRepository["[]string"] == nil },
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.testDescription, func(t *testing.T) {
			result, err := Parse(testCase.input)
			assert.Equal(t, err, testCase.expectedError)
			assert.Equal(t, result.declaredStructs, testCase.expectedStructDeclarations)
			//fmt.Println(result.typesRepository)
			if !testCase.resultPredicate(&result) {
				t.Errorf("Failed to evaluate test predicate")
			}
		})
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
