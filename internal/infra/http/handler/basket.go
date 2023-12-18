package handler

import (
	"errors"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/SahandNoey/Cart-Service/internal/domain/model"
	"github.com/SahandNoey/Cart-Service/internal/domain/repository/basketrepo"
	"github.com/SahandNoey/Cart-Service/internal/infra/http/request"
	"github.com/labstack/echo/v4"
)

const (
	PENDING   string = "PENDING"
	COMPLETED string = "COMPLETED"
)

type BasketH struct {
	repo basketrepo.Repository
}

func NewBasketH(repo basketrepo.Repository) *BasketH {
	return &BasketH{
		repo: repo,
	}
}

func (basketH *BasketH) Get(c echo.Context) error {
	baskets := basketH.repo.Get(c.Request().Context(), basketrepo.GetCommand{})

	if len(baskets) == 0 {
		return echo.ErrNotFound
	}

	return c.JSON(http.StatusOK, baskets)
}

func (basketH *BasketH) GetById(c echo.Context) error {
	var idPtr *uint64
	if id, err := strconv.ParseUint(c.Param("id"), 10, 64); err == nil {
		idPtr = &id
	} else {
		return echo.ErrBadRequest
	}

	basket := basketH.repo.Get(c.Request().Context(), basketrepo.GetCommand{Id: idPtr})

	if len(basket) == 0 {
		return echo.ErrNotFound
	}

	if len(basket) > 1 {
		return echo.ErrInternalServerError
	}

	return c.JSON(http.StatusOK, basket)
}

func (basketH *BasketH) Create(c echo.Context) error {
	var req request.CreateBasket

	if err := c.Bind(&req); err != nil {
		return echo.ErrBadRequest
	}

	if err := req.Validate(); err != nil {
		return echo.ErrBadRequest
	}

	id := rand.Uint64() % 1_000_000
	if err := basketH.repo.Create(c.Request().Context(), model.Basket{
		Id:        id,
		Data:      req.Data,
		State:     PENDING,
		CreatedAt: time.Now(),
	}); err != nil {
		if errors.Is(err, basketrepo.ErrorDuplicateBasketID) {
			return echo.ErrBadRequest
		}
		return echo.ErrInternalServerError
	}
	return c.JSON(http.StatusCreated, id)
}

func (basketH *BasketH) Update(c echo.Context) error {
	var req request.UpdateBasket

	if err := c.Bind(&req); err != nil {
		return echo.ErrBadRequest
	}

	if err := req.Validate(); err != nil {
		return echo.ErrBadRequest
	}

	var idPtr *uint64
	if id, err := strconv.ParseUint(c.Param("id"), 10, 64); err != nil {
		return echo.ErrBadRequest
	} else {
		idPtr = &id
	}

	baskets := basketH.repo.Get(c.Request().Context(), basketrepo.GetCommand{Id: idPtr})
	if len(baskets) == 0 {
		return echo.ErrNotFound
	}
	if len(baskets) > 1 {
		return echo.ErrInternalServerError
	}

	basket := baskets[0]

	if basket.State == COMPLETED {
		return echo.ErrBadRequest
	} else {
		if req.Data != "" {
			basket.Data = req.Data
		}

		if req.State != "" {
			basket.State = req.State
		}
	}

	if err := basketH.repo.Update(c.Request().Context(), model.Basket{
		Id:        basket.Id,
		Data:      basket.Data,
		State:     basket.State,
		UpdatedAt: time.Now(),
		CreatedAt: basket.CreatedAt,
	}); err != nil {
		return echo.ErrInternalServerError
	}

	baskets = basketH.repo.Get(c.Request().Context(), basketrepo.GetCommand{Id: &basket.Id})
	basket = baskets[0]

	return c.JSON(http.StatusOK, basket)
}

func (basketH *BasketH) Delete(c echo.Context) error {
	var idPtr *uint64
	if id, err := strconv.ParseUint(c.Param("id"), 10, 64); err != nil {
		return echo.ErrBadRequest
	} else {
		idPtr = &id
	}

	baskets := basketH.repo.Get(c.Request().Context(), basketrepo.GetCommand{Id: idPtr})
	if len(baskets) == 0 {
		return echo.ErrBadRequest
	}
	if len(baskets) > 1 {
		return echo.ErrInternalServerError
	}

	if err := basketH.repo.Delete(c.Request().Context(), basketrepo.GetCommand{Id: idPtr}); err != nil {
		return echo.ErrInternalServerError
	}

	return c.JSON(http.StatusOK, baskets[0])

}

func (basketH *BasketH) Register(group *echo.Group) {
	group.GET("", basketH.Get)
	group.POST("", basketH.Create)
	group.PATCH(":id", basketH.Update)
	group.GET(":id", basketH.GetById)
	group.DELETE(":id", basketH.Delete)
}
