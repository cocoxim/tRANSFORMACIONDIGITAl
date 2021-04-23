package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/swaggo/swag/example/celler/httputil"
	"github.com/swaggo/swag/example/celler/model"
)

// ShowBottle godoc
//
//	@Summary		Show a bottle
//	@Description	get string by ID
//	@ID				get-string-by-int
//	@Tags			bottles
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"Bottle ID"
//	@Success		200	{object}	model.Bottle
//	@Failure		400	{object}	httputil.HTTPError
//	@Failure		404	{object}	httputil.HTTPError
//