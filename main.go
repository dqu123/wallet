package main

import (
	"github.com/dqu123/wallet-backend/controllers"
	"github.com/dqu123/wallet-backend/database"
	"github.com/dqu123/wallet-backend/models"

	"github.com/gin-gonic/gin"
)

func addControllerRoutes(context *database.Context, router *gin.Engine) {
	controllers.AddWalletControllerRoutes(context, router)
}

func main() {
	router := gin.Default()
	context, err := database.InitDB()
	if err != nil {
		panic(err)
	}
	models.AutoMigrateModels(context)
	addControllerRoutes(context, router)
	router.Run("localhost:8080")
}
