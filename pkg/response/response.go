package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Meta struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Error   interface{} `json:"error"`
}

type APIError struct {
	Code    int    `json:"code,omitempty"`
	Type    string `json:"type,omitempty"`
	Field   string `json:"field,omitempty"`
	Message string `json:"message"`
}

type Paginator struct {
	Total  int64 `json:"total"`
	Limit  int64 `json:"limit"`
	Offset int64 `json:"offset"`
	Link   Link  `json:"links"`
}

type Link struct {
	NextPageUrl string `json:"next_page_url"`
	PrevPageUrl string `json:"prev_page_url"`
}

type MetaPaginator struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Error   interface{} `json:"error"`
	Page    Paginator   `json:"page"`
}

type Single struct {
	Meta Meta        `json:"meta"`
	Data interface{} `json:"data,omitempty"`
}

type Paging struct {
	MetaPaginator MetaPaginator `json:"meta"`
	Data          interface{}   `json:"data,omitempty"`
}

func NotFound(c *gin.Context, message string, error interface{}) {
	c.JSON(http.StatusNotFound, Single{
		Meta: Meta{
			Code:    http.StatusNotFound,
			Message: message,
			Error:   error,
		},
		Data: nil,
	})
	return
}

func DataWithoutMeta(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, data)
	return
}

func BadRequest(c *gin.Context, message string, error interface{}) {
	c.JSON(http.StatusBadRequest, Single{
		Meta: Meta{
			Code:    http.StatusBadRequest,
			Message: message,
			Error:   error,
		},
		Data: nil,
	})
	return
}

func InternalServerError(c *gin.Context, message string, error interface{}) {
	c.JSON(http.StatusInternalServerError, Single{
		Meta: Meta{
			Code:    http.StatusInternalServerError,
			Message: message,
			Error:   error,
		},
		Data: nil,
	})
	return
}

func CustomResponse(c *gin.Context, statusCode int, message string, error interface{}) {
	c.JSON(statusCode, Single{
		Meta: Meta{
			Code:    statusCode,
			Message: message,
			Error:   error,
		},
		Data: nil,
	})
	return
}
