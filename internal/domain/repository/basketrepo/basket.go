package basketrepo

import (
	"context"
	"errors"
	"time"

	"github.com/SahandNoey/Cart-Service/internal/domain/model"
)

var ErrorDuplicateBasketID = errors.New("Given basket ID already exists")

type GetCommand struct {
	Id        *uint64
	CreatedAt *time.Time
	UpdatedAt *time.Time
	Data      *string
	State     *string
}

type Repository interface {
	Get(ctx context.Context, cmd GetCommand) []model.Basket
	Create(ctx context.Context, basket model.Basket) error
	Update(ctx context.Context, basket model.Basket) error
	Delete(ctx context.Context, cmd GetCommand) error
}
