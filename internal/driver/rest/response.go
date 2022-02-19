package rest

import (
	"errors"
	"net/http"
	"time"

	"github.com/go-chi/render"
)

type Response struct {
	StatusCode int         `json:"-"`
	OK         bool        `json:"ok"`
	Data       interface{} `json:"data,omitempty"`
	Err        string      `json:"err,omitempty"`
	Message    string      `json:"msg,omitempty"`
	Timestamp  int64       `json:"ts"`
}

func (rb *Response) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, rb.StatusCode)
	rb.Timestamp = time.Now().Unix()
	return nil
}

func NewSuccessResp(data interface{}) Response {
	return Response{
		StatusCode: http.StatusOK,
		OK:         true,
		Data:       data,
	}
}

func NewErrorResp(err error) *Response {
	var restErr *Error
	if !errors.As(err, &restErr) {
		restErr = NewInternalServerError(err.Error())
	}
	return &Response{
		StatusCode: restErr.StatusCode,
		OK:         false,
		Err:        restErr.Err,
		Message:    restErr.Message,
	}
}
