package rest

import (
	"errors"
	"fmt"
	"net/http"
)

type Error struct {
	StatusCode int
	Err        string
	Message    string
}

func (e *Error) Error() string {
	return fmt.Sprintf("%v - %v - %v", e.StatusCode, e.Err, e.Message)
}

func (e *Error) Is(target error) bool {
	var restErr *Error
	if !errors.As(target, &restErr) {
		return false
	}
	return *e == *restErr
}

func NewInternalServerError(msg string) *Error {
	return &Error{
		StatusCode: http.StatusInternalServerError,
		Err:        "ERR_INTERNAL_ERROR",
		Message:    msg,
	}
}

func NewBadRequestError(msg string) *Error {
	return &Error{
		StatusCode: http.StatusBadRequest,
		Err:        "ERR_BAD_REQUEST",
		Message:    msg,
	}
}

func NewPartnerNotFoundError() *Error {
	return &Error{
		StatusCode: http.StatusNotFound,
		Err:        "ERR_PARTNER_NOT_FOUND",
		Message:    "given `partner_id` is not found",
	}
}

func NewGameNotFoundError() *Error {
	return &Error{
		StatusCode: http.StatusNotFound,
		Err:        "ERR_GAME_NOT_FOUND",
		Message:    "game is not found",
	}
}

func NewBattleNotFoundError() *Error {
	return &Error{
		StatusCode: http.StatusNotFound,
		Err:        "ERR_BATTLE_NOT_FOUND",
		Message:    "battle is not found",
	}
}

func NewInvalidBattleStateError() *Error {
	return &Error{
		StatusCode: http.StatusConflict,
		Err:        "ERR_INVALID_BATTLE_STATE",
		Message:    "invalid battle state",
	}
}
