package model

type exampleGenerator func() string

//DataType is the base interface that all types representing a data type
//should implement.
type DataType interface {
	IsSimple() bool
	IsList() bool
	IsStruct() bool
	Generate() string
}
