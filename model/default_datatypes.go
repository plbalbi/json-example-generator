package model

import (
	"math/rand"
	"strconv"
	"time"
)

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

//TODO: Improve this. Doing it as quick as possible

//DataTypeRepository is the central collection of all DataTypes
type DataTypeRepository map[string]DataType

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

//TODO: Add SimpleDataType constructor
//TODO: Improve random generators

//GetDefaultDataTypeRepository bootrstraps the initial repository, of the most
//common or default DataTypes. From them all others all built.
func GetDefaultDataTypeRepository() DataTypeRepository {
	repository := make(DataTypeRepository)
	randomGenerator := rand.New(rand.NewSource(time.Now().UnixNano()))

	//string
	repository["string"] = &SimpleDataType{
		name: "string",
		generator: func() string {
			length := randomGenerator.Intn(20)
			randString := make([]rune, length)
			for i := range randString {
				randString[i] = letters[randomGenerator.Intn(len(letters))]
			}
			return `"` + string(randString) + `"`
		},
	}
	//int
	repository["int"] = &SimpleDataType{
		name: "int",
		generator: func() string {
			return strconv.Itoa(randomGenerator.Int() % 10000)
		},
	}
	//float64
	repository["float64"] = &SimpleDataType{
		name: "float64",
		generator: func() string {
			return strconv.FormatFloat(randomGenerator.Float64()*1000, 'f', 3, 32)
		},
	}
	//bool
	repository["bool"] = &SimpleDataType{
		name: "bool",
		generator: func() string {
			return []string{"true", "false"}[randomGenerator.Int()%2]
		},
	}
	return repository
}
