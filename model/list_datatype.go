package model

import (
	"bytes"
	"math/rand"
	"time"
)

const maxRandomListSize = 10

//ListDataType represents a list type, containing multiple elements of a certain
//DataType.
type ListDataType struct {
	name      string
	innerType DataType
}

//IsSimple return true if the datatype is a ListDataType.
func (data *ListDataType) IsSimple() bool {
	return false
}

//IsList return true if the datatype is a ListDataType.
func (data *ListDataType) IsList() bool {
	return true
}

//IsStruct return true if the datatype is a ListDataType.
func (data *ListDataType) IsStruct() bool {
	return false
}

//Generate generates a random example of this datatype.
func (data *ListDataType) Generate() string {
	randomGenerator := rand.New(rand.NewSource(time.Now().UnixNano()))
	//Use as list size a random number, mod 10, + 1 to assure that is >= 1
	listSize := randomGenerator.Int()%maxRandomListSize + 1
	var randomListBuffer bytes.Buffer
	randomListBuffer.WriteString("[")
	for i := 1; i <= listSize; i++ {
		randomListBuffer.WriteString(data.innerType.Generate())
		if i < listSize {
			randomListBuffer.WriteString(", ")
		}
	}
	randomListBuffer.WriteString("]")
	return randomListBuffer.String()
}
