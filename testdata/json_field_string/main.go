package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type MyStruct struct {
	ID       int      `json:"id" example:"1" format:"int64"`
	Name     string   `json:"name" example:"poti"`
	Intvar   int      `json:"myint,string"`                            // integer as string
	Boolvar  bool     `json:",string"`                                 // boolean as a string
	TrueBool bool     `json:"truebool,string" example:"true"`          // boolean as a string
	Floatvar float64  `json:",string"`                                 // float as a string
	UUIDs    []string `json:"uuids" type:"arr