package request

import (
	"fmt"

	"github.com/asaskevich/govalidator"
	"github.com/labstack/echo/v4"
)

type Validator interface {
	Validate() error
}

type CreateBasket struct {
	Data string `json:"data,omitempty"  valid:"length(0|2048),optional"`
}

func (cb CreateBasket) Validate() error {
	if _, err := govalidator.ValidateStruct(cb); err != nil {
		return fmt.Errorf("%w: Data Length of basket is too long(more than 2048)", echo.ErrBadRequest)
	}
	return nil
}

type UpdateBasket struct {
	Data  string `json:"data,omitempty" valid:"length(0|2048),optional"`
	State string `json:"state,omitempty" valid:"in(COMPLETED|PENDING),optional"`
}

func (ub UpdateBasket) Validate() error {
	if _, err := govalidator.ValidateStruct(ub); err != nil {
		return fmt.Errorf("%w: Your update request is not valid", echo.ErrBadRequest)
	}
	return nil
}
