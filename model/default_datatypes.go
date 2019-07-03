package model

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

	//string
	repository["string"] = &SimpleDataType{
		name: "string",
		generator: func() string {
			return "\"holis\""
		},
	}
	//int
	repository["int"] = &SimpleDataType{
		name: "int",
		generator: func() string {
			return "42"
		},
	}
	//float64
	repository["float64"] = &SimpleDataType{
		name: "float64",
		generator: func() string {
			return "66.6"
		},
	}
	//bool
	repository["bool"] = &SimpleDataType{
		name: "bool",
		generator: func() string {
			return "true"
		},
	}
	return repository
}
