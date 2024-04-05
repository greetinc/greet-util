package response

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/greetinc/greet-util/s/date"
	"github.com/greetinc/greet-util/s/log"

	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/teris-io/shortid"
)

type errorResponse struct {
	Meta  ResponseModel `json:"meta"`
	Error string        `json:"error"`
}

type Error struct {
	Response     errorResponse `json:"response"`
	Code         int           `json:"code"`
	ErrorMessage error
}

const (
	E_DUPLICATE            = "duplicate"
	E_NOT_FOUND            = "not_found"
	E_UNPROCESSABLE_ENTITY = "unprocessable_entity"
	E_UNAUTHORIZED         = "unauthorized"
	E_BAD_REQUEST          = "bad_request"
	E_SERVER_ERROR         = "server_error"
	E_VERIFIED             = "not_verified"
	E_MAXLIMIT             = "max_limit"
	E_COMPANYEJECT         = "company_eject"
	E_REGISTERMAIL         = "mail_missing"
)

type errorConstant struct {
	Duplicate           Error
	NotFound            Error
	RouteNotFound       Error
	UnprocessableEntity Error
	Unauthorized        Error
	BadRequest          Error
	Validation          Error
	InternalServerError Error
	Unverified          Error
	ExpiredToken        Error
	Suspend             Error
	VerifyPassword      Error
	MaxLimit            Error
	CompanyEject        Error
	RegisterMail        Error
	RecordNotFound      Error
}

var ErrorConstant errorConstant = errorConstant{

	RegisterMail: Error{
		Response: errorResponse{
			Meta: ResponseModel{
				Status:  false,
				Message: "missing '@'",
			},
			Error: E_REGISTERMAIL,
		},
		Code: http.StatusBadRequest,
	},
	CompanyEject: Error{
		Response: errorResponse{
			Meta: ResponseModel{
				Status:  false,
				Message: "Waiting for company data to be approved 1 X 24 hours",
			},
			Error: E_COMPANYEJECT,
		},
		Code: http.StatusBadRequest,
	},
	MaxLimit: Error{
		Response: errorResponse{
			Meta: ResponseModel{
				Status:  false,
				Message: "user has reached the maximum limit for adding data",
			},
			Error: E_MAXLIMIT,
		},
		Code: http.StatusConflict,
	},
	Duplicate: Error{
		Response: errorResponse{
			Meta: ResponseModel{
				Status:  false,
				Message: "Created value already exists",
			},
			Error: E_DUPLICATE,
		},
		Code: http.StatusConflict,
	},
	NotFound: Error{
		Response: errorResponse{
			Meta: ResponseModel{
				Status:  false,
				Message: "Data not found",
			},
			Error: E_NOT_FOUND,
		},
		Code: http.StatusNotFound,
	},
	RouteNotFound: Error{
		Response: errorResponse{
			Meta: ResponseModel{
				Status:  false,
				Message: "Route not found",
			},
			Error: E_NOT_FOUND,
		},
		Code: http.StatusNotFound,
	},
	UnprocessableEntity: Error{
		Response: errorResponse{
			Meta: ResponseModel{
				Status:  false,
				Message: "Unprocessable Invalid parameters or payload",
			},
			Error: E_UNPROCESSABLE_ENTITY,
		},
		Code: http.StatusUnprocessableEntity,
	},
	Unauthorized: Error{
		Response: errorResponse{
			Meta: ResponseModel{
				Status:  false,
				Message: "Unauthorized, please login",
			},
			Error: E_UNAUTHORIZED,
		},
		Code: http.StatusUnauthorized,
	},
	Unverified: Error{
		Response: errorResponse{
			Meta: ResponseModel{
				Status:  false,
				Message: "Account not verified. Please check your email for verification instructions",
			},
			Error: E_VERIFIED,
		},
		Code: http.StatusForbidden,
	},
	RecordNotFound: Error{
		Response: errorResponse{
			Meta: ResponseModel{
				Status:  false,
				Message: "User not found",
			},
			Error: E_BAD_REQUEST,
		},
		Code: http.StatusForbidden,
	},
	ExpiredToken: Error{
		Response: errorResponse{
			Meta: ResponseModel{
				Status:  false,
				Message: "OTP has expired, please resend",
			},
			Error: E_VERIFIED,
		},
		Code: http.StatusForbidden,
	},
	Suspend: Error{
		Response: errorResponse{
			Meta: ResponseModel{
				Status:  false,
				Message: "Your account has been temporarily suspended for 5 minutes following three unsuccessful login attempts.",
			},
			Error: E_VERIFIED,
		},
		Code: http.StatusForbidden,
	},
	VerifyPassword: Error{
		Response: errorResponse{
			Meta: ResponseModel{
				Status:  false,
				Message: "Incorrect email or password. Please try again.",
			},
			Error: E_VERIFIED,
		},
		Code: http.StatusForbidden,
	},
	BadRequest: Error{
		Response: errorResponse{
			Meta: ResponseModel{
				Status:  false,
				Message: "Bad Request",
			},
			Error: E_BAD_REQUEST,
		},
		Code: http.StatusBadRequest,
	},
	Validation: Error{
		Response: errorResponse{
			Meta: ResponseModel{
				Status:  false,
				Message: "Invalid parameters or payload",
			},
			Error: E_BAD_REQUEST,
		},
		Code: http.StatusBadRequest,
	},
	InternalServerError: Error{
		Response: errorResponse{
			Meta: ResponseModel{
				Status:  false,
				Message: "Something bad happened",
			},
			Error: E_SERVER_ERROR,
		},
		Code: http.StatusInternalServerError,
	},
}

func ErrorBuilder(res *Error, message error) *Error {
	res.ErrorMessage = message
	return res
}

func CustomErrorBuilder(code int, err string, message string) *Error {
	return &Error{
		Response: errorResponse{
			Meta: ResponseModel{
				Status:  false,
				Message: message,
			},
			Error: err,
		},
		Code: code,
	}
}

func ErrorResponse(err error) *Error {
	re, ok := err.(*Error)
	if ok {
		return re
	} else {
		return ErrorBuilder(&ErrorConstant.InternalServerError, err)
	}
}

func (e *Error) Error() string {
	return fmt.Sprintf("error code %d", e.Code)
}

func (e *Error) ParseToError() error {
	return e
}

func (e *Error) Send(c echo.Context) error {
	var errorMessage string
	if e.ErrorMessage != nil {
		errorMessage = fmt.Sprintf("%+v", errors.WithStack(e.ErrorMessage))
	}
	logrus.Error(errorMessage)

	body, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		logrus.Warn("error read body, message : ", e.Error())
	}

	bHeader, err := json.Marshal(c.Request().Header)
	if err != nil {
		logrus.Warn("error read header, message : ", e.Error())
	}

	go func() {
		retries := 3
		logError := log.LogError{
			ID:           shortid.MustGenerate(),
			Header:       string(bHeader),
			Body:         string(body),
			URL:          c.Request().URL.Path,
			HttpMethod:   c.Request().Method,
			ErrorMessage: errorMessage,
			Level:        "Error",
			AppName:      os.Getenv("APP"),
			Version:      os.Getenv("VERSION"),
			Env:          os.Getenv("ENV"),
			CreatedAt:    *date.DateTodayLocal(),
		}
		for i := 0; i < retries; i++ {
			err := log.InsertErrorLog(context.Background(), &logError)
			if err == nil {
				break
			}
		}
	}()

	return c.JSON(e.Code, e.Response)
}
