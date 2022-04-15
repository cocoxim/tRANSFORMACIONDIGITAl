package api

import (
	"log"
	"net/http"
	"time"

	"github.com/swaggo/swag/testdata/alias_type/data"
)

/*// @Summary Get time as string
// @Description get time as string
// @ID time-as-string
// @Accept  json
// @Produce  json
// @Success 200 {object} data.StringAlias	"ok"
// @Router /testapi/time-as-string [get]
func GetTimeAsStringAlias(w http.ResponseWriter, r *http.Request) {
	var foo data.StringAlias = "test"
	log.Println(foo)
	//write your code
}*/

/*// @Summary Get time as time
// @Description get time as time
// @ID time-as-time
// @Accept  json
// @Produce  json
// @Success 200 {object} data.DateOnly	"ok"
// @Router /testapi/time-as-time [get]
func GetTimeAsTimeAlias(w http.ResponseWriter, r *http.Request) {
	var foo = data.DateOnly(time.Now())
	log.Println(foo)
	//write your code
}*/

// @