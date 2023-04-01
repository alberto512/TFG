package santander

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func AuthorizeOpenBusiness(code string) (string, error) {

	url := "https://apis-sandbox.bancosantander.es/canales-digitales/sb/v2/authorize/?client_id=bc75ee49-9924-4160-904e-6b246d751e2c&redirect_uri=https://tfg-app.netlify.app&response_type=code"

	payload := strings.NewReader("{\"access\":{\"accounts\":[],\"balances\":[],\"transactions\":[],\"cards_accounts\":[],\"cards_balances\":[],\"cards_transactions\":[]},\"recurringIndicator\":\"true\",\"frequencyPerDay\":\"5\",\"validUntil\":\"2023-09-20\"}")

	req, _ := http.NewRequest("POST", url, payload)

	req.Header.Add("Authorization", code)
	req.Header.Add("content-type", "application/json")
	req.Header.Add("accept", "application/json")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	fmt.Println(res)
	fmt.Println(string(body))

	return "", nil
}