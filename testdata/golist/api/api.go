package api

/*
#include "foo.h"
*/
import "C"
import (
	"fmt"
	"net/http"
)

func PrintInt(i, j int) {
	res := C.add(C.int(i), C.int(j))
	fmt.Println(res)
}

type Foo struct {
	ID        int      `json:"id"`
	Na