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
//	@Produce		