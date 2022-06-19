package web

import (
	"time"
)

type GenericBody[T any] struct {
	Data T
}

type GenericBodyMulti[T a