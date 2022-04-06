package api

import (
	"log"
	"net/http"

	"github.com/swaggo/swag/testdata/alias_import/data"
	"github.com/swaggo/swag/testdata/alias_type/types"
)

// @Summary Get application
// @Description test get application
// @ID get-application
// @Accept  json
// @Produce  json
// @Success 200 {object} data.ApplicationResponse	"ok"
// @Router /testa