package controllers

import (
	"reflect"

	"github.com/gofiber/fiber/v2"
)

func wrapInArray(value interface{}) interface{} {
	// Check if value is already a slice
	v := reflect.ValueOf(value)
	if v.Kind() == reflect.Slice {
		// If value is a slice of string, return it as is
		if reflect.TypeOf(value).Elem().Kind() == reflect.String {
			return value
		}
		// Convert to slice of interface{}
		interfaceSlice := make([]interface{}, v.Len())
		for i := 0; i < v.Len(); i++ {
			interfaceSlice[i] = v.Index(i).Interface()
		}
		return interfaceSlice
	}
	return []interface{}{value}
}

func SuccessResponse(c *fiber.Ctx, statusCode int, data interface{}) error {
	return c.Status(statusCode).JSON(Response{
		Code:    statusCode,
		Success: true,
		Data:    wrapInArray(data),
	})
}

func ErrorResponse(c *fiber.Ctx, statusCode int, errMessage string) error {
	return c.Status(statusCode).JSON(Response{
		Code:    statusCode,
		Success: false,
		Error:   wrapInArray(errMessage),
	})
}
