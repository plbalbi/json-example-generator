package model

import (
	"bytes"
	"fmt"
)

type structFieldMap map[string]string

//StructDataType represents a struct type, which is a map-like structure. And
//also is the main type to be outputed.
type StructDataType struct {
	name   string
	fields structFieldMap
}

//NewStructDataType returns a brand new StructDataType, with the given name
//and no fields in it.
func NewStructDataType(aName string) *StructDataType {
	return &StructDataType{
		name:   aName,
		fields: make(structFieldMap),
	}
}

//GetName shows the datatype's name
func (data *StructDataType) GetName() string {
	return data.name
}

//IsSimple return true if the datatype is a StructDataType.
func (data *StructDataType) IsSimple() bool {
	return false
}

//IsList return true if the datatype is a StructDataType.
func (data *StructDataType) IsList() bool {
	return false
}

//IsStruct return true if the datatype is a StructDataType.
func (data *StructDataType) IsStruct() bool {
	return true
}

func getType(repository DataTypeRepository, typeName string) DataType {
	if typeName[0] != '[' {
		return repository[typeName]
	} else {
		innerType := getType(repository, typeName[2:])
		return &ListDataType{
			name:      typeName,
			innerType: innerType,
		}
	}
}

//Generate generates a random example of this datatype.
func (data *StructDataType) Generate(repository DataTypeRepository) string {
	fmt.Println("generating struct...")
	var randomStructBuffer bytes.Buffer
	randomStructBuffer.WriteString("{")
	printedFieldCounter := 1
	lastFieldNumber := len(data.fields)
	for fieldName, typeName := range data.fields {
		fmt.Println("found field named", fieldName, "of type", typeName)
		fieldType := getType(repository, typeName)
		fmt.Println("type is", fieldType)
		randomStructBuffer.WriteString(fmt.Sprintf("\"%s\": %s", fieldName, fieldType.Generate(repository)))
		if printedFieldCounter < lastFieldNumber {
			randomStructBuffer.WriteString(",")
		}
		printedFieldCounter++
	}
	randomStructBuffer.WriteString("}")
	return randomStructBuffer.String()
}

//TODO: handle errors when:
//	- already a field with name aName
//	- circular datatype definitions

//AddFieldNamed adds a new field to the StructDataType.
func (data *StructDataType) AddFieldNamed(aName string, aTypeName string) error {
	data.fields[aName] = aTypeName
	return nil
}
