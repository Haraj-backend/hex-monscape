package rest

import (
	"errors"
	"net/http"
	"time"

	"github.com/go-chi/render"
)

type RespBody struct {
	StatusCode int         `json:"-"`
	OK         bool        `json:"ok"`
	Data       interface{} `json:"data,omitempty"`
	Err        string      `json:"err,omitempty"`
	Message    string      `json:"msg,omitempty"`
	Timestamp  int64       `json:"ts"`
}

func (rb *RespBody) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, rb.StatusCode)
	rb.Timestamp = time.Now().Unix()
	return nil
}

func NewSuccessResp(data interface{}) *RespBody {
	return &RespBody{
		StatusCode: http.StatusOK,
		OK:         true,
		Data:       data,
	}
}

func NewErrorResp(err error) *RespBody {
	var restErr *Error
	if !errors.As(err, &restErr) {
		restErr = NewInternalServerError(err.Error())
	}
	return &RespBody{
		StatusCode: restErr.StatusCode,
		OK:         false,
		Err:        restErr.Err,
		Message:    restErr.Message,
	}
}
