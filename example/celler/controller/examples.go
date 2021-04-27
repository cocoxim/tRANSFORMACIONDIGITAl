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
//	@Failure		500		{string}	string	"ok"
//	@Router			/examples/calc [get]
func (c *Controller) CalcExample(ctx *gin.Context) {
	val1, err := strconv.Atoi(ctx.Query("val1"))
	if err != nil {
		httputil.NewError(ctx, http.StatusBadRequest, err)
		return
	}
	val2, err := strconv.Atoi(ctx.Query("val2"))
	if err != nil {
		httputil.NewError(ctx, http.StatusBadRequest, err)
		return
	}
	ans := val1 + val2
	ctx.String(http.StatusOK, "%d", ans)
}

// PathParamsExample godoc
//
//	@Summary		path params example
//	@Description	path params
//	@Tags			example
//	@Accept			json
//	@Produce		plain
//	@Param			group_id	path		int		true	"Group ID"
//	@Param			account_id	path		int		true	"Account ID"
//	@Success		200			{string}	string	"answer"
//	@Failure		400			{string}	string	"ok"
//	@Failure		404			{string}	string	"ok"
//	@Failure		500			{string}	string	"ok"
//	@Router			/examples/groups/{group_id}/accounts/{account_id} [get]
func (c *Controller) PathParamsExample(ctx *gin.Context) {
	groupID, err := strconv.Atoi(ctx.Param("group_id"))
	if err != nil {
		httputil.NewError(ctx, http.StatusBadRequest, err)
		return
	}
	accountID, err := strconv.Atoi(ctx.Param("account_id"))
	if err != nil {
		httputil.NewError(ctx, http.StatusBadRequest, err)
		return
	}
	ctx.String(http.StatusOK, "group_id=%d account_id=%d", groupID, accountID)
}

// HeaderExample godoc
//
//	@Summary		custome header example
//	@Description	custome header
//	@Tags			example
//	@Accept			json
//	@Produce		plain
//	@Param			Authorization	header		string	true	"Authentication header"
//	@Success		200				{string}	string	"answer"
//	@Failure		400				{string}	string	"ok"
//	@Failure		404				{string}	string	"ok"
//	@Failure		500				{string}	string	"ok"
//	@Router			/examples/header [get]
func (c *Controller) HeaderExample(ctx *gin.Context) {
	ctx.String(http.StatusOK, ctx.GetHeader("Authorization"))
}

// SecuritiesExample godoc
//
//	@Summary		custome header example
//	@Description	custome header
//	@Tags			example
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string	true	"Authentication header"
//	@Success		200				{string}	string	"answer"
//	@Failure		400				{string}	string	"ok"
//	@Failure		404				{string}	string	"ok"
//	@Failure		500				{string}	string	"ok"
//	@Security		ApiKeyAuth
//	@Security		OAuth2Implicit[admin, write]
//	@Router			/examples/securities [get]
func (c *Controller) SecuritiesExample(ctx *gin.Context) {
}

// AttributeExample godoc
//
//	@Summary		attribute example