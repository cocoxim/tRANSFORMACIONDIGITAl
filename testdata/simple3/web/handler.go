package web

import (
	"time"

	uuid "github.com/gofrs/uuid"
	"github.com/shopspring/decimal"
)

type Pet struct {
	ID       int `example:"1" format:"int64"`
	Category struct {
		ID            int      `example:"1"`
		Name          string   `example:"category_name"`
		PhotoURLs     []string `example:"http://test/image/1.jpg,http://test/image/2.jpg" format:"url"`
		SmallCategory struct {
			ID        int      `example:"1"`
			Name      string   `example:"detail_category_name"`
			PhotoURLs []string `example:"http://test/image/1.jpg,http://test/image/2.jpg"`
		}
	}
	Name      string   `example:"poti"`
	PhotoURLs []string `example:"http://test/image/1.jpg,http://test/image/2.jpg"`
	Tags      []Tag
	Pets      *[]Pet2
	Pets2     []*Pet2