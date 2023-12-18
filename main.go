package main

import (
	"fmt"
	"log"

	"github.com/SahandNoey/Cart-Service/internal/domain/repository/basketrepo"
	"github.com/SahandNoey/Cart-Service/internal/infra/http/handler"
	"github.com/SahandNoey/Cart-Service/internal/infra/repository/basketsql"
	"github.com/labstack/echo/v4"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	dbConfig := "host=localhost user=postgres password=Sahand9923087 dbname=CartServiceDB port=5432 sslmode=disable TimeZone=Asia/Tehran"
	db, err := gorm.Open(postgres.Open(dbConfig), &gorm.Config{})
	if err != nil {
		fmt.Printf("error connecting to database: %v", err)
	}

	if err = db.AutoMigrate(new(basketsql.BasketDTO)); err != nil {
		fmt.Printf("failed to automigrate: %v", err)
	}

	app := echo.New()
	var repo basketrepo.Repository = basketsql.New(db)

	h := handler.NewBasketH(repo)
	h.Register(app.Group("basket/"))

	if err := app.Start("0.0.0.0:1377"); err != nil {
		log.Fatalf("server failed to start %v", err)
	}
}
