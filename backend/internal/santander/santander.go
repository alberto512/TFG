package santander

import (
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

const redirect_uri = "https://tfg-app.netlify.app/"
const endpoint = "https://apis-sandbox.bancosantander.es/canales-digitales/sb/v2/";

const tokenEndpoint = endpoint + "token";

type ResponseTokenEndpoint struct {
	AccessToken		string	`json:"access_token"`
	TokenType  		string 	`json:"token_type"`
	ExpiresIn   	int 	`json:"expires_in"`
	RefreshToken	string 	`json:"refresh_token"`
	AuthToken 		string 	`json:"auth_token"`
}

func GetTokenWithCode(code string) (string, error) {
	log.Printf("Get token with code")

	// Create body
	body := url.Values{}
	body.Set("grant_type", "authorization_code")
	body.Set("redirect_uri", redirect_uri)
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
	req.Header.Add("Authorization", b64.StdEncoding.EncodeToString([]byte(os.Getenv("SANTANDER_ID")+":"+os.Getenv("SANTANDER_ID"))))
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

	fmt.Println("ResponseTokenEndpoint:", response)
	fmt.Println("AccessToken:", response.AccessToken)
	fmt.Println("TokenType:", response.TokenType)
	fmt.Println("ExpiresIn:", response.ExpiresIn)
	fmt.Println("RefreshToken:", response.RefreshToken)
	fmt.Println("AuthToken:", response.AuthToken)

	return "", nil
}

func GetTokenWithRefresh(refresh string) (string, error) {
	log.Printf("Get token with code")

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
	req.Header.Add("Authorization", b64.StdEncoding.EncodeToString([]byte(os.Getenv("SANTANDER_ID")+":"+os.Getenv("SANTANDER_ID"))))
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

	fmt.Println("ResponseTokenEndpoint:", response)
	fmt.Println("AccessToken:", response.AccessToken)
	fmt.Println("TokenType:", response.TokenType)
	fmt.Println("ExpiresIn:", response.ExpiresIn)
	fmt.Println("RefreshToken:", response.RefreshToken)
	fmt.Println("AuthToken:", response.AuthToken)

	return "", nil
}