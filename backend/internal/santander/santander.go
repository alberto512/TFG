package santander

import (
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"tfg/internal/mongo"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const redirectUri = "https://tfg-frontend-production.up.railway.app/santanderLogin/"
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

type ResponseAccount struct {
	Iban		string	`json:"iban"`
	Name  		string 	`json:"name"`
	Currency	string 	`json:"currency"`
}

type ResponseAccountsEndpoint struct {
	AccountList	[]ResponseAccount	`json:"accountList"`
	RequestId	string				`json:"requestId"`
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

func GetToken(userId string) (string, error) {
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

func GetAccount(accessToken string, iban string) (string, error) {
	log.Printf("Get account of user")
	

	// Create the request
	req, err := http.NewRequest("GET", accountsEndpoint + "/" + iban + "?withBalance=true", nil)
	if err != nil {
		log.Printf("Error: Create request")
        return "", err
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
        return "", err
	}
	if res.StatusCode != http.StatusOK {
		log.Printf("Error: Response %d", res.StatusCode)
        return "", fmt.Errorf("error %d", res.StatusCode)
	}
	defer res.Body.Close()

	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("client: could not read response body: %s\n", err)
		os.Exit(1)
	}
	fmt.Printf("client: response body: %s\n", resBody)

	var response [2][]byte
	response[0] = resBody

	// Create body
	body := []byte(`{
		"movement": "BOTH",
		"date_to": "2023-04-12",
		"date_from": "2023-01-15",
		"amount_to": 100000000,
		"amount_from": 0,
		"order": "A"
	}`)

	payload := strings.NewReader("{\"movement\":\"BOTH\",\"date_to\":\"2023-04-12\",\"date_from\":\"2023-01-15\",\"amount_to\":100000000,\"amount_from\":0,\"order\":\"A\"}")

	// Create the request
	req, err = http.NewRequest("POST", movementsEndpoint, payload)
	if err != nil {
		log.Printf("Error: Create request")
        return "", err
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
        return "", err
	}
	if res.StatusCode != http.StatusOK {
		log.Printf("Error: Response %d", res.StatusCode)
		resBody, err = ioutil.ReadAll(res.Body)
		if err != nil {
			fmt.Printf("client: could not read response body: %s\n", err)
			os.Exit(1)
		}
		fmt.Printf("client: response body: %s\n", resBody)
        return "", fmt.Errorf("error %d", res.StatusCode)
	}
	defer res.Body.Close()

	resBody, err = ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("client: could not read response body: %s\n", err)
		os.Exit(1)
	}
	fmt.Printf("client: response body: %s\n", resBody)

	response[1] = resBody

	// Decode the response
	/*
	response := &ResponseAccountsEndpoint{}
	derr := json.NewDecoder(res.Body).Decode(response)
	if derr != nil {
		log.Printf("Error: Decoding response")
        return "", err
	}

	log.Printf("Accounts %v", response)
	out, err := json.Marshal(response)
    if err != nil {
        panic(err)
    }
	*/

	out, err := json.Marshal(response)
    if err != nil {
        panic(err)
    }

	return string(out), nil
}

func GetAccounts(accessToken string) (string, error) {
	log.Printf("Get accounts of user")

	// Create the request
	req, err := http.NewRequest("GET", accountsEndpoint, nil)
	if err != nil {
		log.Printf("Error: Create request")
        return "", err
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
        return "", err
	}
	if res.StatusCode != http.StatusOK {
		log.Printf("Error: Response %d", res.StatusCode)
        return "", fmt.Errorf("error %d", res.StatusCode)
	}
	defer res.Body.Close()

	// Decode the response
	response := &ResponseAccountsEndpoint{}
	derr := json.NewDecoder(res.Body).Decode(response)
	if derr != nil {
		log.Printf("Error: Decoding response")
        return "", err
	}

	log.Printf("Accounts %v", response)

	var resFinal [2]string

	for index, element := range response.AccountList {
		resFinal[index], _ = GetAccount(accessToken, element.Iban)
	}

	out, err := json.Marshal(resFinal)
    if err != nil {
        panic(err)
    }

	return string(out), nil
}