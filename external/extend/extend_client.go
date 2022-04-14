package extend

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/dqu123/wallet-backend/database"
	"github.com/dqu123/wallet-backend/logger"
	"github.com/dqu123/wallet-backend/models"
)

const (
	ProviderExtend           = "extend"
	ExtendAPIHost            = "api.paywithextend.com"
	HeaderApplicationVersion = "application/vnd.paywithextend.v2021-03-12+json"

	TokenExpiration = 10 * time.Minute
)

type ExtendClient struct {
	Context *database.Context
}

type ExtendTransactionsResponseBody struct {
	Transactions []ExtendTransaction `json:"transactions"`
}
type ExtendTransaction struct {
	ID                     string `json:"id"`
	MCC                    string `json:"mcc"`
	MCCGroup               string `json:"mccGroup"`
	MCCDescription         string `json:"mccDescription"`
	MerchantName           string `json:"merchantName"`
	AuthBillingAmountCents int    `json:"authBillingAmountCents"`
	AuthBillingCurrency    string `json:"authBillingCurrency"`
	AuthedAt               string `json:"authedAt"`
}

func (ec *ExtendClient) getAccessToken() (string, error) {
	db := ec.Context.DB
	accessToken := &models.AccessToken{}
	tx := db.Last(accessToken)
	if tx.Error != nil || accessToken == nil || accessToken.ExpirationUTC == nil || time.Now().After(*accessToken.ExpirationUTC) {
		postLoginResponse, err := postAccessToken()
		if err != nil {
			return "", err
		}
		expirationUTC := time.Now().Add(TokenExpiration)
		newAccessToken := &models.AccessToken{
			TokenValue:    postLoginResponse.AccessToken,
			ExpirationUTC: &expirationUTC,
		}
		db.Create(newAccessToken)
		return postLoginResponse.AccessToken, nil
	}
	return accessToken.TokenValue, nil
}

func (ec *ExtendClient) GetTransactions(providerCardID string) ([]models.Transaction, error) {
	accessToken, err := ec.getAccessToken()
	// postLoginResponse, err := postAccessToken()
	if err != nil {
		logger.LogError("Error getting access token", err)
		return nil, err
	}
	// accessToken := postLoginResponse.AccessToken
	fmt.Println("accessToken", accessToken)
	fmt.Println("ProviderCardID", providerCardID)
	req, err := http.NewRequest(
		"GET",
		fmt.Sprintf("https://%s/virtualcards/%s/transactions", ExtendAPIHost, providerCardID),
		nil,
	)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", HeaderApplicationVersion)
	headerAuthorization := fmt.Sprintf("Bearer %s", accessToken)
	req.Header.Set("Authorization", headerAuthorization)
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	switch res.StatusCode {
	case http.StatusUnauthorized:
		return nil, fmt.Errorf("issue accessing %s", ExtendAPIHost)
	}
	resBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	extendResBody := &ExtendTransactionsResponseBody{}
	err = json.Unmarshal(resBytes, extendResBody)
	if err != nil {
		return nil, err
	}

	extendTransactions := extendResBody.Transactions
	normalizedTransactions := []models.Transaction{}
	for _, extendTransaction := range extendTransactions {
		externalTransaction := models.ExternalTransaction{
			ProviderID:             extendTransaction.ID,
			MCC:                    extendTransaction.MCC,
			MCCGroup:               extendTransaction.MCCGroup,
			MCCDescription:         extendTransaction.MCCDescription,
			MerchantName:           extendTransaction.MerchantName,
			AuthBillingAmountCents: extendTransaction.AuthBillingAmountCents,
			AuthBillingCurrency:    extendTransaction.AuthBillingCurrency,
			AuthedUTC:              parseExtendTime(extendTransaction.AuthedAt),
		}
		normalizedTransaction := models.Transaction{
			ProviderName:        ProviderExtend,
			ExternalTransaction: externalTransaction,
		}
		normalizedTransactions = append(normalizedTransactions, normalizedTransaction)
	}
	return normalizedTransactions, nil
}

type ExtendVirtualCardsResponseBody struct {
	ExtendCards []ExtendCard `json:"virtualCards"`
}
type ExtendCard struct {
	ID                 string `json:"id" gorm:"column:provider_id;unique"`
	Status             string `json:"status"`
	DisplayName        string `json:"displayName"`
	Expires            string `json:"expires"`
	Currency           string `json:"currency"`
	LimitCents         int    `json:"limitCents"`
	BalanceCents       int    `json:"balanceCents"`
	LifetimeSpentCents int    `json:"lifetimeSpentCents"`
	LastFour           string `json:"last4"`
}

func (ec *ExtendClient) GetVirtualCards() ([]models.VirtualCard, error) {
	accessToken, err := ec.getAccessToken()
	// postLoginResponse, err := postAccessToken()
	if err != nil {
		fmt.Println("Error getting access token", err)
		return nil, err
	}
	// accessToken := postLoginResponse.AccessToken
	fmt.Println("accessToken", accessToken)
	req, err := http.NewRequest(
		"GET",
		fmt.Sprintf("https://%s/virtualcards", ExtendAPIHost),
		nil,
	)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", HeaderApplicationVersion)
	headerAuthorization := fmt.Sprintf("Bearer %s", accessToken)
	req.Header.Set("Authorization", headerAuthorization)
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	switch res.StatusCode {
	case http.StatusUnauthorized:
		return nil, fmt.Errorf("issue accessing %s", ExtendAPIHost)
	}
	resBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	extendResBody := &ExtendVirtualCardsResponseBody{}
	err = json.Unmarshal(resBytes, extendResBody)
	if err != nil {
		return nil, err
	}

	extendCards := extendResBody.ExtendCards
	return normalizeCards(extendCards), nil
}

func normalizeCards(extendCards []ExtendCard) []models.VirtualCard {
	normalizedCards := []models.VirtualCard{}
	for _, extendCard := range extendCards {
		externalCard := models.ExternalVirtualCard{
			ProviderID:         extendCard.ID,
			Status:             extendCard.Status,
			DisplayName:        extendCard.DisplayName,
			ExpirationUTC:      parseExtendTime(extendCard.Expires),
			Currency:           extendCard.Currency,
			LimitCents:         extendCard.LimitCents,
			BalanceCents:       extendCard.BalanceCents,
			LifetimeSpentCents: extendCard.LifetimeSpentCents,
			LastFour:           extendCard.LastFour,
		}
		normalizedCard := models.VirtualCard{
			ProviderName:        ProviderExtend,
			ExternalVirtualCard: externalCard,
		}
		normalizedCards = append(normalizedCards, normalizedCard)
	}
	return normalizedCards
}

type PostLoginRequestBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
type PostLoginResponseBody struct {
	AccessToken string `json:"token"`
}

func postAccessToken() (*PostLoginResponseBody, error) {
	fmt.Println("Posting access token")
	// TODO (WLT-0): Add user management
	email := os.Getenv("USER_EMAIL")
	password := os.Getenv("USER_PASSWORD")
	requestBody := PostLoginRequestBody{
		Email:    email,
		Password: password,
	}

	requestBodyBytes, err := json.Marshal(requestBody)
	if err != nil {
		return nil, err
	}

	// Get token
	req, err := http.NewRequest(
		"POST",
		fmt.Sprintf("https://%s/signin", ExtendAPIHost),
		bytes.NewBuffer(requestBodyBytes),
	)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", HeaderApplicationVersion)
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		// Crashes server if network is down
		log.Fatal(err)
		return nil, err
	}

	// read response data
	resBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	responseBody := &PostLoginResponseBody{}
	err = json.Unmarshal(resBytes, responseBody)
	if err != nil {
		return nil, err
	}
	return responseBody, nil
}
