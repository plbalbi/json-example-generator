package model

import (
	"bytes"
	"fmt"
)

type structFieldMap map[string]DataType

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

//Generate generates a random example of this datatype.
func (data *StructDataType) Generate() string {
	var randomStructBuffer bytes.Buffer
	randomStructBuffer.WriteString("{")
	printedFieldCounter := 1
	lastFieldNumber := len(data.fields)
	for fieldName, fieldType := range data.fields {
		randomStructBuffer.WriteString(fmt.Sprintf("\"%s\": %s", fieldName, fieldType.Generate()))
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
func (data *StructDataType) AddFieldNamed(aName string, aType DataType) error {
	data.fields[aName] = aType
	return nil
}
