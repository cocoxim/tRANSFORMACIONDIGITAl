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
		PhotoUrls     []string `json:"photo_urls" example:"http://test/image/1.jpg,http://test/image/2.jpg" format:"url"`
		SmallCategory struct {
			ID        int      `json:"id" example:"1"`
			Name      string   `json:"name" example:"detail_category_name" binding:"required" minLength:"4" maxLength:"16"`
			PhotoUrls []string `json:"photo_urls" example:"http://test/image/1.jpg,http://test/image/2.jpg"`
		} `json:"small_category"`
	} `json:"category"`
	Name              string            `json:"name" example:"poti" binding:"required"`
	PhotoUrls         []string          `json:"photo_urls" example:"http://test/image/1.jpg,http://test/image/2.jpg" binding:"required"`
	Tags              []Tag             `json:"tags"`
	Pets              *[]Pet2           `json:"pets"`
	Pets2             []*Pet2           `json:"pets2"`
	Status            string            `json:"status" enums:"healthy,ill"`
	Price             float32           `json:"price" example:"3.25" minimum:"1.0" maximum:"1000" multipleOf:"0.01"`
	IsAlive           bool              `json:"is_alive" example:"true" default:"true"`
	Data              interface{}       `json:"data"`
	Hidden            string            `json:"-"`
	UUID              uuid.UUID         `json:"uuid"`
	Decima