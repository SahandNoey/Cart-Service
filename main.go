package main

import (
	"log"

	"github.com/SahandNoey/Cart-Service/internal/domain/repository/basketrepo"
	"github.com/SahandNoey/Cart-Service/internal/infra/http/handler"
	"github.com/SahandNoey/Cart-Service/internal/infra/repository/basketmem"
	"github.com/labstack/echo/v4"
)

func main() {
	app := echo.New()
	var repo basketrepo.Repository = basketmem.New()

	h := handler.NewBasketH(repo)
	h.Register(app.Group("basket/"))

	if err := app.Start("0.0.0.0:3000"); err != nil {
		log.Fatalf("server failed to start %v", err)
	}
}
