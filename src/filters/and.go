package filters

import "strings"

type And struct {
	Filters []Filter
}

func NewAnd(filters ...Filter) And {
	return And{Filters: filters}
}

func (a And) GetSQL() string {
	var sqlParts []string
	for _, filter := range a.Filters {
		sqlParts = append(sqlParts, filter.GetSQL())
	}
	return strings.Join(sqlParts, " AND ")
}

func (a And) GetParams() []interface{} {
	var params []interface{}
	for _, filter := range a.Filters {
		for _, param := range filter.GetParams() {
			params = append(params, param)
		}
	}
	return params
}
