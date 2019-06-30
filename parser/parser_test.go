package parser

import "testing"

func TestParser(t *testing.T) {
	result, err := Parse("type perro struct {}")
	if err != nil {
		t.Errorf("Test failed with error: %s", err.Error())
	}
	if result.structsCount != 1 {
		t.Errorf("Expected %d structs to be parsed, but found %d", 1, result.structsCount)
	}
}
