package response

import (
	dto "aseprayana-skeleton-go/dto/transaction"
	"net/http"

	"github.com/labstack/echo/v4"
)

type successConstant struct {
	OK Success
}

var SuccessConstant successConstant = successConstant{
	OK: Success{
		Response: successResponse{
			Code:    http.StatusOK,
			Status:  true,
			Message: "Request successfully proceed",
			Data:    nil,
		},
	},
}

type successResponse struct {
	Status  bool        `json:"status"`
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type Success struct {
	Response successResponse `json:"response"`
	Code     int             `json:"code"`
}

func SuccessBuilder(res *Success, data interface{}) *Success {
	res.Response.Data = data
	return res
}

func CustomSuccessBuilder(code int, data interface{}, message string, info *dto.PaginationResponse) *Success {
	return &Success{
		Response: successResponse{
			Status:  true,
			Message: message,
			Data:    data,
		},
	}

}

func SuccessResponse(data interface{}) *Success {
	return SuccessBuilder(&SuccessConstant.OK, data)
}

func (s *Success) Send(c echo.Context) error {
	return c.JSON(s.Code, s.Response)
}
