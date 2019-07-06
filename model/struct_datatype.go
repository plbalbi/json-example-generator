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

//GenerateWithIndentationPrefix generates a random example of this datatype, with a provided indetation level.
func (data *StructDataType) GenerateWithIndentationPrefix(repository DataTypeRepository, indentationPrefix string) string {
	var randomStructBuffer bytes.Buffer
	randomStructBuffer.WriteString("{\n")
	printedFieldCounter := 1
	lastFieldNumber := len(data.fields)
	for fieldName, typeName := range data.fields {
		fieldType := getType(repository, typeName)
		if fieldType.IsStruct() {
			fieldTypeAsStruct := fieldType.(*StructDataType)
			randomStructBuffer.WriteString(fmt.Sprintf("%s\t\"%s\": %s",
				indentationPrefix,
				fieldName,
				fieldTypeAsStruct.GenerateWithIndentationPrefix(
					repository,
					fmt.Sprintf("\t%s", indentationPrefix))))
		} else {
			randomStructBuffer.WriteString(fmt.Sprintf("%s\t\"%s\": %s", indentationPrefix, fieldName, fieldType.Generate(repository)))
		}
		if printedFieldCounter < lastFieldNumber {
			randomStructBuffer.WriteString(",\n")
		}
		printedFieldCounter++
	}
	randomStructBuffer.WriteString(fmt.Sprintf("\n%s}", indentationPrefix))
	return randomStructBuffer.String()
}

//Generate generates a random example of this datatype.
func (data *StructDataType) Generate(repository DataTypeRepository) string {
	return data.GenerateWithIndentationPrefix(repository, "")
}

//AddFieldNamed adds a new field to the StructDataType.
func (data *StructDataType) AddFieldNamed(aName string, aTypeName string) error {
	data.fields[aName] = aTypeName
	return nil
}
