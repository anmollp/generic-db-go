package filters

import (
	"fmt"
	"strings"
)

type In struct {
	Column string
	Values []interface{}
}

func (i In) GetSQL() string {
	if len(i.Values) == 0 {
		// "col IN ()" is not valid, so we return "false" when the filter is empty.
		return "false"
	}
	placeholders := strings.Repeat("?, ", len(i.Values)-1) + "?"
	return fmt.Sprintf("%s IN (%s)", i.Column, placeholders)
}

func (i In) GetParams() []interface{} {
	var params []interface{}
	for _, param := range i.Values {
		params = append(params, param)
	}
	return params
}

type NotIn struct {
	Column string
	Values []interface{}
}

func (n NotIn) GetSQL() string {
	if len(n.Values) == 0 {
		return "true"
	}
	placeholders := strings.Repeat("?, ", len(n.Values)-1) + "?"
	return fmt.Sprintf("%s NOT IN (%s)", n.Column, placeholders)
}

func (n NotIn) GetParams() []interface{} {
	var params []interface{}
	for _, param := range n.Values {
		params = append(params, param)
	}
	return params
}

type MultiColumnIn struct {
	Columns []string
	Values  [][]interface{}
}

func (m MultiColumnIn) GetSQL() string {
	if len(m.Values) == 0 {
		// "col IN ()" is not valid, so we return "false" when the filter is empty.
		return "false"
	}
	columnList := "(" + strings.Join(m.Columns, ", ") + ")"
	valuePlaceholders := "(" + strings.Repeat("?, ", len(m.Columns)-1) + "?)"
	allPlaceholders := strings.Repeat(valuePlaceholders+", ", len(m.Values)-1) + valuePlaceholders
	return fmt.Sprintf("%s IN (%s)", columnList, allPlaceholders)
}

func (m MultiColumnIn) GetParams() []interface{} {
	var flatParams []interface{}
	for _, valueGroup := range m.Values {
		for _, val := range valueGroup {
			flatParams = append(flatParams, val)
		}
	}
	return flatParams
}
