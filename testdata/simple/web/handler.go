package web

import (
	"time"

	"github.com/gofrs/uuid"
	"github.com/shopspring/decimal"
	"github.com/swaggo/swag/testdata/simple/cross"
)

type Pet struct {
	ID       int `json:"id" example:"1" format:"int64" readonly:"true"`
	Category struct {
		ID            int      `json:"id" example:"1"`
		Name          string   `json:"name" example:"category_name"`
		PhotoUrls     []string `json:"ph