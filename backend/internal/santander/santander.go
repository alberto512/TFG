package santander

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func AuthorizeOpenBusiness(code string) (string, error) {

	url := "https://apis-sandbox.bancosantander.es/canales-digitales/sb/v2/authorize/?client_id=bc75ee49-9924-4160-904e-6b246d751e2c&redirect_uri=https://tfg-app.netlify.app&response_type=code"

	payload := strings.NewReader("{\"access\":{\"accounts\":[{\"currency\":\"EUR\",\"msisdn\":\"…\",\"bban\":\"20385778983000760236\",\"iban\":\"ES1111111111111111111\"}],\"balances\":[{\"currency\":\"EUR\",\"msisdn\":\"…\",\"bban\":\"20385778983000760236\",\"iban\":\"ES1111111111111111111\"}],\"transactions\":[{\"currency\":\"EUR\",\"msisdn\":\"…\",\"bban\":\"20385778983000760236\",\"iban\":\"ES1111111111111111111\"}],\"allPsd2\":\"allAccounts\",\"availableAccounts\":\"allAccounts\",\"availableAccountsWithBalances\":\"6376496895865296\",\"cards_accounts\":[{\"resourceId\":\"732692424631\"}],\"cards_balances\":[{\"resourceId\":\"732692424631\"}],\"cards_transactions\":[{\"resourceId\":\"732692424631\"}]},\"recurringIndicator\":\"true\",\"frequencyPerDay\":\"5\",\"validUntil\":\"2019-09-20\",\"scopes\":\"[\\\"myprofile.read\\\"]\"}")

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