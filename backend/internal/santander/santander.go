package santander

import (
	"bytes"
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"tfg/internal/accounts"
	"tfg/internal/mongo"
	"tfg/internal/transactions"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const redirectUri = "https://tfg-frontend.up.railway.app/santanderLogin/"
const endpoint = "https://apis-sandbox.bancosantander.es/canales-digitales/sb/v2/";

const tokenEndpoint = endpoint + "token";
const accountsEndpoint = endpoint + "accounts";
const movementsEndpoint = endpoint + "movements";

type ResponseTokenEndpoint struct {
	AccessToken		string	`json:"access_token"`
	TokenType  		string 	`json:"token_type"`
	ExpiresIn   	int 	`json:"expires_in"`
	RefreshToken	string 	`json:"refresh_token"`
}

type UserToken struct {
    ID       		string		`bson:"_id,omitempty"`
	UserID      	string		`bson:"userId,omitempty"`
    AccessToken 	string 		`bson:"accessToken,omitempty"`
    Expires 		time.Time	`bson:"expires,omitempty"`
	RefreshToken	string 		`bson:"refreshToken,omitempty"`
}

type ResponseAmountBalance struct {
	Currency	string	`json:"currency"`
	Content		string	`json:"content"`
}

type ResponseTransactions struct {
	BookingDate string					`json:"bookingDate"`
	ValueDate   string					`json:"valueDate"`
	Description	string					`json:"description"`
	Amount      ResponseAmountBalance	`json:"amount"`
	Balance     ResponseAmountBalance 	`json:"balance"`
}

type ResponseAccountWithTransactions struct {
	Iban			string					`json:"iban"`
	Transactions	[]ResponseTransactions	`json:"transactions"`
}

type ResponseTransactionsEndpoint struct {
	Account	ResponseAccountWithTransactions	`json:"account"`
}

type ResponseBalance struct {
	Amount					string	`json:"amount"`
	Currency				string 	`json:"currency"`
	CreditDebitIndicator	string	`json:"creditDebitIndicator"`
}

type ResponseAccountExpanded struct {
	Iban		string				`json:"iban"`
	Name  		string				`json:"name"`
	Currency	string				`json:"currency"`
	Balance		ResponseBalance 	`json:"balance"`
}

type ResponseAccountEndpoint struct {
	Account		ResponseAccountExpanded	`json:"account"`
	RequestId	string					`json:"requestId"`
}

type ResponseAccount struct {
	Iban		string	`json:"iban"`
	Name  		string 	`json:"name"`
	Currency	string 	`json:"currency"`
}

type ResponseAccountsEndpoint struct {
	AccountList	[]ResponseAccount	`json:"accountList"`
	RequestId	string				`json:"requestId"`
}

type RequestBodyAccount struct {
	Movement	string	`json:"movement"`
	DateTo		string	`json:"date_to"`
	DateFrom	string	`json:"date_from"`
	AmountTo	int		`json:"amount_to"`
	AmountFrom	int		`json:"amount_from"`
	Order		string	`json:"order"`
}

func saveToken(userId string, token *ResponseTokenEndpoint) (error) {
	log.Printf("Save token in database")

	// String id to ObjectId
	id, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		log.Printf("Error: Convert string to id")
        return err
	}

	// Get date of expiration of token
	date := time.Now().UnixMilli() + (int64(token.ExpiresIn) * 1000)

	// Filter to get token by userId
	filter := bson.D{
		{Key: "userId", Value: id},
	}

	// Update fields
	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "userId", Value: id},
			{Key: "accessToken", Value: token.AccessToken},
			{Key: "expires", Value: date},
			{Key: "refreshToken", Value: token.RefreshToken},
			},
		},
	}

	// Set option to create if it doesn't exist
	opts := options.Update().SetUpsert(true)

	// Execute insert
	_, err = mongo.UpdateOne("santanderTokens", filter, update, opts)
	if err != nil {
		log.Printf("Error: Save operation in db")
		return err
	}

	return nil
}

func getTokenWithRefresh(userId string, refresh string) (string, error) {
	log.Printf("Get token with refresh")

	// Create body
	body := url.Values{}
	body.Set("grant_type", "refresh_token")
	body.Set("refresh_token", refresh)

	// Encode the body
	encodedBody := body.Encode()

	// Create the request
	req, err := http.NewRequest("POST", tokenEndpoint, strings.NewReader(encodedBody))
	if err != nil {
		log.Printf("Error: Create request")
        return "", err
	}
	
	// Add all the headers
	req.Header.Add("X-IBM-Client-Id", os.Getenv("SANTANDER_ID"))
	req.Header.Add("X-IBM-Client-Secret", os.Getenv("SANTANDER_SECRET"))
	req.Header.Add("Authorization", b64.StdEncoding.EncodeToString([]byte(os.Getenv("SANTANDER_ID")+":"+os.Getenv("SANTANDER_SECRET"))))
	req.Header.Add("content-type", "application/x-www-form-urlencoded")
	req.Header.Add("accept", "application/json")

	// Make the request
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("Error: Make request")
        return "", err
	}
	if res.StatusCode != http.StatusOK {
		log.Printf("Error: Response %d", res.StatusCode)
        return "", fmt.Errorf("error %d", res.StatusCode)
	}
	defer res.Body.Close()

	// Decode the response
	response := &ResponseTokenEndpoint{}
	derr := json.NewDecoder(res.Body).Decode(response)
	if derr != nil {
		log.Printf("Error: Decoding response")
        return "", err
	}

	// Save tokens in database
	err = saveToken(userId, response)
	if err != nil {
		log.Printf("Error: Save token")
        return "", err
	}

	return response.AccessToken, nil
}

func getToken(userId string) (string, error) {
	var userToken *UserToken

	log.Printf("Get token by userId")

	// String id to ObjectId
	id, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		log.Printf("Error: Convert string to id")
        return "", err
	}

	// Query to get token by userId
	query := bson.D{
		{Key: "userId", Value: id},
	}

	// Empty filter
    filter := bson.D{}

	// Execute query
    cursor, err := mongo.Query("santanderTokens", query, filter)
    if err != nil {
        log.Printf("Error: Get tokens in db")
        return "", err
    }

    // Decode query
    cursor.Next(mongo.GetCtx())

    if err = cursor.Decode(&userToken); err != nil {
        log.Printf("Error: Decoding user token %s", err)
        return "", err
    }

	if time.Now().After(userToken.Expires) {
		return getTokenWithRefresh(userId, userToken.RefreshToken)
	}

	return userToken.AccessToken, nil
}

func refreshAccount(accessToken string, iban string, userId string) (error) {
	log.Printf("Refresh account")
	
	// Create the request
	req, err := http.NewRequest("GET", accountsEndpoint + "/" + iban + "?withBalance=true", nil)
	if err != nil {
		log.Printf("Error: Create request")
        return err
	}
	
	// Add all the headers
	req.Header.Add("Authorization", "Bearer " + accessToken)
	req.Header.Add("X-IBM-Client-Id", os.Getenv("SANTANDER_ID"))
	req.Header.Add("accept", "application/json")
	req.Header.Add("psu_active", "1")

	// Make the request
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("Error: Make request")
        return err
	}
	if res.StatusCode != http.StatusOK {
		log.Printf("Error: Response %d", res.StatusCode)
        return fmt.Errorf("error %d", res.StatusCode)
	}
	defer res.Body.Close()

	// Decode the response
	response := &ResponseAccountEndpoint{}
	err = json.NewDecoder(res.Body).Decode(response)
	if err != nil {
		log.Printf("Error: Decoding response")
        return err
	}

	// Save account in database
	var account accounts.Account
	account.Iban = response.Account.Iban
	account.Name = response.Account.Name
	account.Currency = response.Account.Currency
	amount, err := strconv.ParseFloat(response.Account.Balance.Amount, 64)
	if err != nil {
		log.Printf("Error: Parsing amount")
        return err
	}
	account.Amount = amount
	account.Bank = "Santander"
	account.UserID = userId

	updateDate := time.Date(1970, time.January, 1, 0, 0, 0, 0, time.UTC)

	// Check if account exists. Create it if not or update it if it does
	err = account.GetAccountByIban(userId)
	if err != nil {
		err = account.Create()
		if err != nil {
			log.Printf("Error: Creating account")
			return err
		}
	} else {
		updateDate = account.UpdateDate
		err = account.Update(userId, nil, &response.Account.Name, &response.Account.Currency, &amount, nil)
		if err != nil {
			log.Printf("Error: Updating account")
			return err
		}
	}

	dateTo := time.Now().Format(time.DateOnly)
	dateFrom := updateDate.Format(time.DateOnly)

	days := time.Since(updateDate).Hours() / 24

	if days > 90 {
		dateFrom = time.Now().AddDate(0, 0, -90).Format(time.DateOnly)
	}

	if dateFrom == dateTo {
		return nil
	}

	// Create body
	body := RequestBodyAccount{
		Movement: "BOTH",
		DateTo: dateTo,
		DateFrom: dateFrom,
		AmountTo: 100000000,
		AmountFrom: 0,
		Order: "A",
	}

	fmt.Println(body)

	parsedBody, err := json.Marshal(body)
	if err != nil {
		log.Printf("Error: Encoding body")
		return err
	}

	fmt.Println(parsedBody)
	
	// Create the request
	req, err = http.NewRequest("POST", movementsEndpoint + "/" + iban, bytes.NewBuffer(parsedBody))
	if err != nil {
		log.Printf("Error: Create request")
        return err
	}
	
	// Add all the headers
	req.Header.Add("Authorization", "Bearer " + accessToken)
	req.Header.Add("X-IBM-Client-Id", os.Getenv("SANTANDER_ID"))
	req.Header.Add("content-type", "application/json")
	req.Header.Add("accept", "application/json")
	req.Header.Add("psu_active", "1")

	// Make the request
	res, err = http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("Error: Make request")
        return err
	}
	if res.StatusCode != http.StatusOK {
		log.Printf("Error: Response %d", res.StatusCode)
        return fmt.Errorf("error %d", res.StatusCode)
	}
	defer res.Body.Close()

	// Decode the response
	responseTransactions := &ResponseTransactionsEndpoint{}
	err = json.NewDecoder(res.Body).Decode(responseTransactions)
	if err != nil {
		log.Printf("Error: Decoding response")
        return err
	}

	// Save transactions in database
	for _, element := range responseTransactions.Account.Transactions {
		var transaction transactions.Transaction
		
		transaction.Description = element.Description
		transaction.Date, err = time.Parse(time.DateOnly, element.ValueDate)
		if err != nil {
			log.Printf("Error: Parsing date")
			return err
		}
		transaction.Amount, err = strconv.ParseFloat(element.Amount.Content, 64)
		if err != nil {
			log.Printf("Error: Parsing amount")
			return err
		}
		transaction.Category = "Default"
		transaction.AccountID = account.ID

		err = transaction.Create()
		if err != nil {
			log.Printf("Error: Creating transaction")
			return err
		}
	}

	return nil
}

func GetTokenWithCode(userId string, code string) (string, error) {
	log.Printf("Get token with code")

	// Create body
	body := url.Values{}
	body.Set("grant_type", "authorization_code")
	body.Set("redirect_uri", redirectUri)
	body.Set("code", code)

	// Encode the body
	encodedBody := body.Encode()

	// Create the request
	req, err := http.NewRequest("POST", tokenEndpoint, strings.NewReader(encodedBody))
	if err != nil {
		log.Printf("Error: Create request")
        return "", err
	}
	
	// Add all the headers
	req.Header.Add("X-IBM-Client-Id", os.Getenv("SANTANDER_ID"))
	req.Header.Add("X-IBM-Client-Secret", os.Getenv("SANTANDER_SECRET"))
	req.Header.Add("Authorization", b64.StdEncoding.EncodeToString([]byte(os.Getenv("SANTANDER_ID")+":"+os.Getenv("SANTANDER_SECRET"))))
	req.Header.Add("content-type", "application/x-www-form-urlencoded")
	req.Header.Add("accept", "application/json")

	// Make the request
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("Error: Make request")
        return "", err
	}
	if res.StatusCode != http.StatusOK {
		log.Printf("Error: Response %d", res.StatusCode)
        return "", fmt.Errorf("error %d", res.StatusCode)
	}
	defer res.Body.Close()

	// Decode the response
	response := &ResponseTokenEndpoint{}
	derr := json.NewDecoder(res.Body).Decode(response)
	if derr != nil {
		log.Printf("Error: Decoding response")
        return "", err
	}

	// Save tokens in database
	err = saveToken(userId, response)
	if err != nil {
		log.Printf("Error: Save token")
        return "", err
	}

	return response.AccessToken, nil
}

func RefreshData(userId string) (error) {
	log.Printf("Refresh data of user")

	accessToken, err := getToken(userId)
	if err != nil {
		log.Printf("Error: Get user token")
        return err
	}

	// Create the request
	req, err := http.NewRequest("GET", accountsEndpoint, nil)
	if err != nil {
		log.Printf("Error: Create request")
        return err
	}
	
	// Add all the headers
	req.Header.Add("Authorization", "Bearer " + accessToken)
	req.Header.Add("X-IBM-Client-Id", os.Getenv("SANTANDER_ID"))
	req.Header.Add("accept", "application/json")
	req.Header.Add("psu_active", "1")

	// Make the request
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("Error: Make request")
        return err
	}
	if res.StatusCode != http.StatusOK {
		log.Printf("Error: Response %d", res.StatusCode)
        return fmt.Errorf("error %d", res.StatusCode)
	}
	defer res.Body.Close()

	// Decode the response
	response := &ResponseAccountsEndpoint{}
	err = json.NewDecoder(res.Body).Decode(response)
	if err != nil {
		log.Printf("Error: Decoding response")
        return err
	}

	// Refresh each account from the response
	for _, element := range response.AccountList {
		err := refreshAccount(accessToken, element.Iban, userId)
		if err != nil {
			log.Printf("Error: Refreshing account")
        	return err
		}
	}

	return nil
}