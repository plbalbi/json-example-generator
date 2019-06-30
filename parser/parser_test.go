package parser

import "testing"

func TestParserCountsASingleStruct(t *testing.T) {
	result, err := Parse("type perro")
	if err != nil {
		t.Errorf("Test failed with error: %s", err.Error())
	}
	if result.structsCount != 1 {
		t.Errorf("Expected %d structs to be parsed, but found %d", 1, result.structsCount)
	}
}

func TestParserCountsTwoStructs(t *testing.T) {
	result, err := Parse("type perro type gato")
	if err != nil {
		t.Errorf("Test failed with error: %s", err.Error())
	}
	if result.structsCount != 2 {
		t.Errorf("Expected %d structs to be parsed, but found %d", 2, result.structsCount)
	}
}
