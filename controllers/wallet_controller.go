package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm/clause"

	"example.com/wallet-backend/database"
	"example.com/wallet-backend/external/extend"
	"example.com/wallet-backend/logger"
	"example.com/wallet-backend/models"
)

type WalletController struct {
	Context      *database.Context
	ExtendClient extend.ExtendClient
}

type GetCardsResponseBody struct {
	VirtualCards []models.VirtualCard `json:"virtualCards"`
}

func (wc WalletController) GetCards(c *gin.Context) {
	db := wc.Context.DB
	virtualCards := &[]models.VirtualCard{}
	tx := db.Find(virtualCards)
	if tx.Error != nil {
		logger.LogError("db.Find(virtualCards)", tx.Error)
	}
	res := &GetCardsResponseBody{
		VirtualCards: *virtualCards,
	}
	c.JSON(http.StatusOK, res)
}

type GetTransactionsResponseBody struct {
	Transactions []models.Transaction `json:"transactions"`
}

func (wc WalletController) GetTransactions(c *gin.Context) {
	db := wc.Context.DB
	transactions := &[]models.Transaction{}
	tx := db.Find(transactions)
	if tx.Error != nil {
		logger.LogError("db.Find(virtualCards)", tx.Error)
	}
	res := &GetTransactionsResponseBody{
		Transactions: *transactions,
	}
	c.JSON(http.StatusOK, res)
}

type IngestDataResponseBody struct {
	VirtualCards []models.VirtualCard
	Transactions []models.Transaction
}

// Ingests Virtual Cards and all
func (wc WalletController) Ingest(c *gin.Context) {
	db := wc.Context.DB
	extendClient := wc.ExtendClient
	fmt.Println("Ingesting")
	virtualCards, err := extendClient.GetVirtualCards()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// TODO (WLT-1): handle DB failure
	db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "provider_id"}},
		UpdateAll: true,
	}).Create(virtualCards)

	transactions, err := extendClient.GetTransactions(virtualCards[0].ProviderID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "provider_id"}},
		UpdateAll: true,
	}).Create(transactions)

	res := &IngestDataResponseBody{
		VirtualCards: virtualCards,
		Transactions: transactions,
	}
	c.JSON(http.StatusOK, res)
}

func AddWalletControllerRoutes(context *database.Context, router *gin.Engine) {
	WalletController := WalletController{
		Context: context,
		ExtendClient: extend.ExtendClient{
			Context: context,
		},
	}
	router.GET("/cards", WalletController.GetCards)
	router.GET("/transactions", WalletController.GetTransactions)

	router.POST("/ingest", WalletController.Ingest)
}
