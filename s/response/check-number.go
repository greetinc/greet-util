package response

import (
	"fmt"

	"github.com/labstack/echo/v4"
)

// func IsNumber(data string) error {
// 	_, err := strconv.Atoi(data)
// 	if err != nil {
// 		return fmt.Errorf("path parameter hanya menerima angka")
// 	}

// 	return nil
// }

func IsNumber(c echo.Context, paramName string) (string, error) {
	value := c.Param(paramName)

	if len(value) == 0 {
		return "", fmt.Errorf("%s cannot be empty", paramName)
	}

	return value, nil
}
