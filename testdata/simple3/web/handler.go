package web

import (
	"time"

	uuid "github.com/gofrs/uuid"
	"github.com/shopspring/decimal"
)

type Pet struct {
	ID       int `example:"