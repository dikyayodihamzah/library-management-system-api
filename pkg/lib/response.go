package lib

import (
	"encoding/json"
	"strconv"

	"github.com/dikyayodihamzah/library-management-api/pkg/model"
	"github.com/dikyayodihamzah/library-management-api/pkg/utils"
	"github.com/gofiber/fiber/v2"
)

// OK send http 200 response
func OK(c *fiber.Ctx, result ...interface{}) error {
	response := model.Response{
		Code:    200,
		Status:  true,
		Message: "Success",
	}
	if len(result) > 0 {
		response.Data = result[0]
	}
	if utils.PrintLog {
		response.Log = AddLog()
	}
	return c.Status(200).JSON(response)
}

// Created send http 201 response
func Created(c *fiber.Ctx, result ...interface{}) error {
	response := model.Response{
		Code:    201,
		Status:  true,
		Message: "Created",
	}
	if len(result) > 0 {
		response.Data = result[0]
	}
	if utils.PrintLog {
		response.Log = AddLog()
	}
	return c.Status(201).JSON(response)
}

func Page[T any](c *fiber.Ctx, total int, result []T) error {
	page := 1
	if i, _ := strconv.Atoi(c.Query("page")); i > 0 {
		page = i
	}

	count := len(result)
	response := model.Response{
		Code:    200,
		Message: "success",
		Status:  true,
		Page:    &page,
		Count:   &count,
		Total:   &total,
		Data:    result,
	}
	if utils.PrintLog {
		response.Log = AddLog()
	}
	return c.Status(200).JSON(response)
}

func AddLog() *string {
	by, _ := json.MarshalIndent(utils.Log, "", "   ")
	utils.Log = []interface{}{}
	return Strptr(string(by))
}
