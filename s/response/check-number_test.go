package response

import (
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

// MockEchoContext is a mock implementation of echo.Context for testing purposes.
type MockEchoContext struct {
	ParamFunc func(name string) string
	echo.Context
}

// Param implements the Param method of echo.Context.
func (m *MockEchoContext) Param(name string) string {
	if m.ParamFunc != nil {
		return m.ParamFunc(name)
	}
	return ""
}

func TestIsNumber_Success(t *testing.T) {
	mockContext := &MockEchoContext{
		ParamFunc: func(name string) string {
			return "123"
		},
	}

	id, err := IsNumber(mockContext, "id")

	assert.Nil(t, err)
	assert.Equal(t, uint(123), id)
}

func TestIsNumber_Failure_InvalidFormat(t *testing.T) {
	mockContext := &MockEchoContext{
		ParamFunc: func(name string) string {
			return "12dwq3"
		},
	}

	id, err := IsNumber(mockContext, "id")

	assert.NotNil(t, err)
	assert.Equal(t, uint(0), id)
	assert.EqualError(t, err, "invalid id format: strconv.ParseUint: parsing \"12dwq3\": invalid syntax")
}

func TestIsNumber_Failure_EmptyID(t *testing.T) {
	mockContext := &MockEchoContext{
		ParamFunc: func(name string) string {
			return ""
		},
	}

	id, err := IsNumber(mockContext, "id")

	assert.NotNil(t, err)
	assert.Equal(t, uint(0), id)
	assert.EqualError(t, err, "invalid id format: strconv.ParseUint: parsing \"\": invalid syntax")
}
