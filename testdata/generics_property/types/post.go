package types

type SubField1[T any, T2 any] struct {
	SubValue1 T
	SubValue2 T2
}

type Field[T any] struct {
	Value  T
	Value2 *