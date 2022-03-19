package main

import (
	"example.com/wallet-backend/controllers"
	"example.com/wallet-backend/database"
	"example.com/wallet-backend/models"

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
