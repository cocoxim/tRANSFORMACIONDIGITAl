package api

import (
	myv1 "github.com/swaggo/swag/testdata/generics_package_alias/internal/path1/v1"
	myv2 "github.com/swaggo/swag/testdata/generics_package_alias/internal/path2/v1"
)

// @Summary Create movie
// @Description Create a new movie production
// @Accept  json
// @Produce  json
// @Success 200 {object} myv2.ListResult[myv2.ProductDto] ""
// @Router /api03 [post]
func CreateMovie03() {
	_ = myv2.ListResult[myv2.ProductDto]{}
}

// @Summary Create movie
// @Description Create a new movie production
// @Accept  json
// @Produce  json
// @Success 200 {object} myv2.RenamedListR