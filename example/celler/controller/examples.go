package controller

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/swaggo/swag/example/celler/httputil"
)

// PingExample godoc
//
//	@Summary		ping example
//	@Description	do ping
//	@Tags			example
//	@Accept			json
//	@Produce		plain
//	@Success		200	{string}	string	"pong"
//	@Failure		400	{string}	string	"ok"
//	@Failure		404	{string}	string	"ok"
//	@Failure		500	{string}	string	"ok"
//	@Router			/examples/ping [get]
func (c *Controller) PingExample(ctx *gin.Context) {
	ctx.String(http.StatusOK, "pong")
}

// CalcExample godoc
//
//	@Summary		calc example
//	@Description	plus
//	@Tags			example
//	@Accept			json
//	@Produce		plain
//	@Param			val1	query		int		true	"used for calc"
//	@Param			val2	query		int		true	"used for calc"
//	@Success		200		{integer}	string	"answer"
//	@Failure		400		{string}	string	"ok"
//	@Failure		404		{string}	string	"ok"
//	@