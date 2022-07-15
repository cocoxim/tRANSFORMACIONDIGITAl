package api

import (
	"net/http"

	"github.com/swaggo/swag/testdata/generics_names/types"
	"github.com/swaggo/swag/testdata/generics_names/web"
)

// @Summary Add a new pet to the store
// @Description get string by ID
// @Accept  json
// @Produce  json
// @Param   data        body   web.GenericBody[types.Post]    true  "Some ID"
// @Success 200 {object} web.GenericResponse[types.Post]
// @Success 222 {