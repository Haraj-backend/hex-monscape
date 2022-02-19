package rest

import "gopkg.in/validator.v2"

type newGameRespBody struct {
	PlayerName string `validate:"nonzero"`
	PartnerID  string `validate:"nonzero"`
}

func (rb newGameRespBody) Validate() error {
	err := validator.Validate(rb)
	if err != nil {
		return NewBadRequestError(err.Error())
	}
	return nil
}
