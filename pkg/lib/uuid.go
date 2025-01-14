package lib

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func GenID() *uuid.UUID {
	u, _ := uuid.NewRandom()
	return &u
}

func IsEmptyUUID(u uuid.UUID) bool {
	return u == uuid.Nil
}

func StringID(u *uuid.UUID) string {
	if u == nil {
		return uuid.Nil.String()
	}
	return u.String()
}

func CompareID(u, x *uuid.UUID) bool {
	if u == nil && x == nil {
		return true
	}

	if u == nil || x == nil {
		return false
	}

	return u.String() == x.String()
}

func ParamsID(c *fiber.Ctx, param string) *uuid.UUID {
	id, err := uuid.Parse(c.Params(param))
	if nil != err {
		return nil
	}

	return &id
}

func QueryID(c *fiber.Ctx, query string) (u []*uuid.UUID) {
	for _, q := range strings.Split(c.Query(query), ",") {
		if id, err := uuid.Parse(q); err == nil {
			u = append(u, Pointer(id))
		}
	}
	return
}
