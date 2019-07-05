package model

import (
	"math/rand"
	"strconv"
	"time"
)

//DataTypeRepository is the central collection of all DataTypes
type DataTypeRepository map[string]DataType
type seedProvider func() int64

const testRandomSeed = 42

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

const stringTypeName = "string"
const intTypeName = "int"
const floatTypeName = "float64"
const boolTypeName = "bool"

//Helpers. Maybe the repository should be abstracted

//CountStructDataTypes counts the struct datatypes registered in a repository
func CountStructDataTypes(repository DataTypeRepository) int {
	count := 0
	for _, datatype := range repository {
		if datatype.IsStruct() {
			count++
		}
	}
	return count
}

func newRandomStringGenerator(randomGenerator *rand.Rand) exampleGenerator {
	return func() string {
		length := randomGenerator.Intn(20)
		randString := make([]rune, length)
		for i := range randString {
			randString[i] = letters[randomGenerator.Intn(len(letters))]
		}
		return `"` + string(randString) + `"`
	}
}

func newRandomIntGenerator(randomGenerator *rand.Rand) exampleGenerator {
	return func() string {
		return strconv.Itoa(randomGenerator.Int() % 10000)
	}
}

func newRandomFloatGenerator(randomGenerator *rand.Rand) exampleGenerator {
	return func() string {
		return strconv.FormatFloat(randomGenerator.Float64()*1000, 'f', 3, 32)
	}
}

func newRandomBoolGenerator(randomGenerator *rand.Rand) exampleGenerator {
	return func() string {
		return []string{"true", "false"}[randomGenerator.Int()%2]
	}
}

//TODO: Add SimpleDataType constructor

func generateRepositoryWithRandomSeed(provider seedProvider) DataTypeRepository {
	repository := make(DataTypeRepository)
	randomGenerator := rand.New(rand.NewSource(provider()))

	//string
	repository[stringTypeName] = &SimpleDataType{
		name:      stringTypeName,
		generator: newRandomStringGenerator(randomGenerator),
	}
	//int
	repository[intTypeName] = &SimpleDataType{
		name:      intTypeName,
		generator: newRandomIntGenerator(randomGenerator),
	}
	//float64
	repository[floatTypeName] = &SimpleDataType{
		name:      floatTypeName,
		generator: newRandomFloatGenerator(randomGenerator),
	}
	//bool
	repository[boolTypeName] = &SimpleDataType{
		name:      boolTypeName,
		generator: newRandomBoolGenerator(randomGenerator),
	}
	return repository
}

//GetDefaultDataTypeRepository bootstraps the initial repository, of the most
//common or default DataTypes. From them all others all built.
func GetDefaultDataTypeRepository() DataTypeRepository {
	return generateRepositoryWithRandomSeed(func() int64 { return time.Now().UnixNano() })
}

//GetTestDataTypeRepository gives a new test DataTypeRepository, which initial datatypes
//use a fixed initial seed to give deterministic results each time.
func GetTestDataTypeRepository() DataTypeRepository {
	return generateRepositoryWithRandomSeed(func() int64 { return testRandomSeed })
}
