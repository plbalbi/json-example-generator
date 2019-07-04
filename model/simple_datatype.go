package model

//SimpleDataType represents a simple type, such as int, string, boolean, etc.
type SimpleDataType struct {
	name      string
	generator exampleGenerator
}

//GetName shows the datatype's name
func (data *SimpleDataType) GetName() string {
	return data.name
}

//IsSimple return true if the datatype is a SimpleDataType.
func (data *SimpleDataType) IsSimple() bool {
	return true
}

//IsList return true if the datatype is a SimpleDataType.
func (data *SimpleDataType) IsList() bool {
	return false
}

//IsStruct return true if the datatype is a SimpleDataType.
func (data *SimpleDataType) IsStruct() bool {
	return false
}

//Generate generates a random example of this datatype.
func (data *SimpleDataType) Generate(repository DataTypeRepository) string {
	return data.generator()
}
