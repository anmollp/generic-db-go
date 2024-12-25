package filters

type Filter interface {
	GetSQL() string
	GetParams() []interface{}
}
