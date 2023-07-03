package model

type Query struct {
	FieldName string
	Boost     float64
	Terms     []string
}
