package model

//TODO: Improve this. Doing it as quick as possible

//DataTypeRepository is the central collection of all DataTypes
type DataTypeRepository map[string]DataType

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
			return "holis"
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
