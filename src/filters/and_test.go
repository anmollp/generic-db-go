package filters

import (
	"reflect"
	"testing"
)

func TestAndFilter_GetSQL(t *testing.T) {
	eqFilter := Equals{Column: "status", Value: "active"}
	gtFilter := NewGreaterEquals("age", 18)

	andFilter := NewAnd(eqFilter, gtFilter)

	expectedSQL := "status = ? AND age >= ?"
	actualSQL := andFilter.GetSQL()

	if actualSQL != expectedSQL {
		t.Errorf("GetSQL() failed. Expected: %s, Got: %s", expectedSQL, actualSQL)
	}
}

func TestAndFilter_GetParams(t *testing.T) {
	eqFilter := Equals{Column: "status", Value: "active"}
	leFilter := NewLessEquals("weight", 200)

	andFilter := NewAnd(eqFilter, leFilter)

	expectedParams := []interface{}{"active", 200}
	actualParams := andFilter.GetParams()

	if !reflect.DeepEqual(actualParams, expectedParams) {
		t.Errorf("GetParams() failed. Expected: %v, Got: %v", expectedParams, actualParams)
	}
}

func TestAndFilter_EmptyFilters(t *testing.T) {
	andFilter := NewAnd()

	expectedSQL := ""
	actualSQL := andFilter.GetSQL()

	if actualSQL != expectedSQL {
		t.Errorf("GetSQL() for empty filters failed. Expected: %s, Got: %s", expectedSQL, actualSQL)
	}

	var expectedParams []interface{}
	actualParams := andFilter.GetParams()

	if !reflect.DeepEqual(actualParams, expectedParams) {
		t.Errorf("GetParams() for empty filters failed. Expected: %v, Got: %v", expectedParams, actualParams)
	}
}
