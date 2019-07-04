package model

import "testing"

func TestEmptyStructIsWellGenerated(t *testing.T) {
	expected := "{}"
	emptyStructType := NewStructDataType("emptyStructType")
	generated := emptyStructType.Generate(nil)
	if expected != generated {
		t.Errorf("Failed to generate type random example. Expected: %s, but got: %s", expected, generated)
	}
}
