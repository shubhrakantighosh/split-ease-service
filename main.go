package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"main/config"
	initilizer "main/init"
	"main/internal/controller"
	opostgres "main/pkg/db/postgres"
	"main/router"
)

func main() {

	ctx := context.TODO()
	config.InitConfig()
	initilizer.Initialize(ctx)

	c := controller.Wire(ctx, opostgres.GetCluster().DbCluster)

	app := gin.New()
	app.GET("/", c.Ge)

	router.PublicRoutes(ctx, app)

	port := viper.GetString("server.port")
	if err := app.Run(port); err != nil {
		panic("failed to start server: " + err.Error())
	}
}
