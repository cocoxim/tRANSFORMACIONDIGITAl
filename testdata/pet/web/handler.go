package web

type Pet struct {
	ID       int `json:"id" example:"1"`
	Category struct {
		ID            int      `json:"id" example:"1"`
		Name          string   `json:"name" example:"category_name"`
		PhotoUrls     []string `json:"photo_urls" example:"http://test/image/1.jpg,http://test/image/2.jpg"`
		SmallCategory struct {
			ID        int      `json:"id" example:"1"`
			Name      string   `json:"name" example:"detail_category_name"`
			PhotoUrls []string `json:"photo_urls" example:"http://test/image/1.jpg,http://tes