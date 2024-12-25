package filters

import "fmt"

type Equals struct {
	Column string
	Value  interface{}
}

func (e Equals) GetSQL() string {
	return fmt.Sprintf("%s = ?", e.Column)
}

func (e Equals) GetParams() []interface{} {
	return []interface{}{e.Value}
}

type GreaterEquals struct {
	Equals
}

func NewGreaterEquals(column string, value interface{}) GreaterEquals {
	return GreaterEquals{
		Equals: Equals{Column: column, Value: value},
	}
}

func (ge GreaterEquals) GetSQL() string {
	return fmt.Sprintf("%s >= ?", ge.Column)
}

type LessEquals struct {
	Equals
}

func NewLessEquals(column string, value interface{}) LessEquals {
	return LessEquals{
		Equals: Equals{Column: column, Value: value},
	}
}

func (le LessEquals) GetSQL() string {
	return fmt.Sprintf("%s <= ?", le.Column)
}
