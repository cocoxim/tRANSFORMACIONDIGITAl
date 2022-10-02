package web

import "github.com/swaggo/swag/testdata/generics_property/types"

type PostSelector func(selector func())

type Filter interface {
	~func(selector func())
}

type query[T any, F Filter] interface {
	Where(ps ...F) T
}

type Pager[T query[T, F], F Filter] struct {
	Rows   uint8   `j