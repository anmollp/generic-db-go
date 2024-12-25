package filters

import (
	"reflect"
	"testing"
)

func TestEqualsFilter_GetSQL(t *testing.T) {
	eqFilter := Equals{Column: "age", Value: 18}
	expectedSQL := "age = ?"
	actualSQL := eqFilter.GetSQL()
	if expectedSQL != actualSQL {
		t.Errorf("GetSQL() failed. Expected: %s, Got: %s", expectedSQL, actualSQL)
	}
}

func TestEqualsFilter_GetParams(t *testing.T) {
	eqFilter := Equals{Column: "age", Value: 18}
	expectedParams := []interface{}{18}
	actualParams := eqFilter.GetParams()
	if !reflect.DeepEqual(expectedParams, actualParams) {
		t.Errorf("GetParams() failed. Expected: %v, Got: %v", expectedParams, actualParams)
	}
}

func TestGreaterEqualsFilter_GetSQL(t *testing.T) {
	eqFilter := NewGreaterEquals("age", 18)
	expectedSQL := "age >= ?"
	actualSQL := eqFilter.GetSQL()
	if expectedSQL != actualSQL {
		t.Errorf("GetSQL() failed. Expected: %s, Got: %s", expectedSQL, actualSQL)
	}
}

func TestGreaterEqualsFilter_GetParams(t *testing.T) {
	eqFilter := NewGreaterEquals("age", 18)
	expectedParams := []interface{}{18}
	actualParams := eqFilter.GetParams()
	if !reflect.DeepEqual(expectedParams, actualParams) {
		t.Errorf("GetParams() failed. Expected: %v, Got: %v", expectedParams, actualParams)
	}
}

func TestLessEqualsFilter_GetSQL(t *testing.T) {
	eqFilter := NewLessEquals("age", 18)
	expectedSQL := "age <= ?"
	actualSQL := eqFilter.GetSQL()
	if expectedSQL != actualSQL {
		t.Errorf("GetSQL() failed. Expected: %s, Got: %s", expectedSQL, actualSQL)
	}
}

func TestLessEqualsFilter_GetParams(t *testing.T) {
	eqFilter := NewLessEquals("age", 18)
	expectedParams := []interface{}{18}
	actualParams := eqFilter.GetParams()
	if !reflect.DeepEqual(expectedParams, actualParams) {
		t.Errorf("GetParams() failed. Expected: %v, Got: %v", expectedParams, actualParams)
	}
}
