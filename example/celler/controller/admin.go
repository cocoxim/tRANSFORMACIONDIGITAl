package controller

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/swaggo/swag/example/celler/httputil"
	"github.com/swaggo/swag/example/celler/model"
)

// Auth godoc
//
//	@Summary		Auth admin
//	@Description	get admin info
//	@Tags			accounts,admin
//	@Accept			json
//	@Produce		json
//	@Success		200	