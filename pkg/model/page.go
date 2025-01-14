package model

import (
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/spf13/cast"
)

type QueryParam struct {
	Page   int    `query:"page" json:"page" example:"1"`           // current page, start from 1
	Limit  int    `query:"limit" json:"limit" example:"10"`        // size per page, if 0 mean all
	Search string `query:"search" json:"search" example:"keyword"` // search keyword
	Sort   string `query:"sort" json:"sort" example:"-created_at"` // add dash in beginning mean desc
}

func GenID() *uuid.UUID {
	u, _ := uuid.NewRandom()
	return &u
}

// TimeNow func
func TimeNow() time.Time {
	t := time.Now()
	return t
}

func StrToIntSlice(s string) []int {
	slc := strings.Split(s, ",")
	var j []int
	for i := range slc {
		j = append(j, cast.ToInt(slc[i]))

	}

	return j
}

func StrToSlice(s string) []string {
	return strings.Split(s, ",")
}
