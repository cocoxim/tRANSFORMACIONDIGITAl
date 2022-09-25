package v1

type ProductDto struct {
	Name1 string `json:"name1"`
}

type ListResult[T any] struct {
	Items1 []T `json:"items1,omitempty"`
}

type RenamedProductDto struct {