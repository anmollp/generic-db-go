package filters

import (
	"reflect"
	"testing"
)

func TestInFilter_GetSQL(t *testing.T) {
	inFilter := In{Column: "forward", Values: []interface{}{9, 10, 11}}
	expectedSQL := "forward IN (?, ?, ?)"
	actualSQL := inFilter.GetSQL()
	if expectedSQL != actualSQL {
		t.Errorf("GetSQL() failed. Expected: %s, Got: %s", expectedSQL, actualSQL)
	}
}

func TestInFilter_GetParams(t *testing.T) {
	inFilter := In{Column: "forward", Values: []interface{}{9, 10, 11}}
	expectedParams := []interface{}{9, 10, 11}
	actualParams := inFilter.GetParams()
	if !reflect.DeepEqual(expectedParams, actualParams) {
		t.Errorf("GetParams() failed. Expected: %v, Got: %v", expectedParams, actualParams)
	}
}

func TestInFilter_EmptyValues(t *testing.T) {
	inFilter := In{Column: "forward", Values: []interface{}{}}
	expectedSQL := "false"
	actualSQL := inFilter.GetSQL()
	if expectedSQL != actualSQL {
		t.Errorf("GetSQL() failed. Expected: %s, Got: %s", expectedSQL, actualSQL)
	}
}

func TestNotInFilter_GetSQL(t *testing.T) {
	notInFilter := NotIn{Column: "forward", Values: []interface{}{4, 5, 6}}
	expectedSQL := "forward NOT IN (?, ?, ?)"
	actualSQL := notInFilter.GetSQL()
	if expectedSQL != actualSQL {
		t.Errorf("GetSQL() failed. Expected: %s, Got: %s", expectedSQL, actualSQL)
	}
}

func TestNotInFilter_GetParams(t *testing.T) {
	notInFilter := NotIn{Column: "forward", Values: []interface{}{4, 5, 6}}
	expectedParams := []interface{}{4, 5, 6}
	actualParams := notInFilter.GetParams()
	if !reflect.DeepEqual(expectedParams, actualParams) {
		t.Errorf("GetParams() failed. Expected: %v, Got: %v", expectedParams, actualParams)
	}
}

func TestNotInFilter_EmptyValues(t *testing.T) {
	notInFilter := NotIn{Column: "forward", Values: []interface{}{}}
	expectedSQL := "true"
	actualSQL := notInFilter.GetSQL()
	if expectedSQL != actualSQL {
		t.Errorf("GetSQL() failed. Expected: %s, Got: %s", expectedSQL, actualSQL)
	}
}

func TestMultiColumnInFilter_GetSQL(t *testing.T) {
	multiInFilter := MultiColumnIn{Columns: []string{"forward", "speed"}, Values: [][]interface{}{
		{9, 33},
		{10, 34},
		{11, 35},
	}}
	expectedSQL := "(forward, speed) IN ((?, ?), (?, ?), (?, ?))"
	actualSQL := multiInFilter.GetSQL()
	if expectedSQL != actualSQL {
		t.Errorf("GetSQL() failed. Expected: %s, Got: %s", expectedSQL, actualSQL)
	}
}

func TestMultiColumnInFilter_GetParams(t *testing.T) {
	multiInFilter := MultiColumnIn{Columns: []string{"forward", "speed"}, Values: [][]interface{}{
		{9, 33},
		{10, 34},
		{11, 35},
	}}
	expectedParams := []interface{}{9, 33, 10, 34, 11, 35}
	actualParams := multiInFilter.GetParams()
	if !reflect.DeepEqual(expectedParams, actualParams) {
		t.Errorf("GetParams() failed. Expected: %v, Got: %v", expectedParams, actualParams)
	}
}

func TestMultiColumnInFilter_EmptyValues(t *testing.T) {
	multiInFilter := MultiColumnIn{Columns: []string{"forward", "speed"}, Values: [][]interface{}{}}
	expectedSQL := "false"
	actualSQL := multiInFilter.GetSQL()
	if expectedSQL != actualSQL {
		t.Errorf("GetSQL() failed. Expected: %s, Got: %s", expectedSQL, actualSQL)
	}
}
