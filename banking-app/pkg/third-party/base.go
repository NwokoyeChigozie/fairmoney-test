package third_party

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type ThirdPartyPackageResponse struct {
	AccountId string  `json:"account_id"`
	Reference string  `json:"reference"`
	Amount    float64 `json:"amount"`
}

type ThirdPartyPkg interface {
	CreateTransaction(accountId, reference string, amount float64) (*ThirdPartyPackageResponse, error)
	GetTransaction(reference string) (*ThirdPartyPackageResponse, error)
}

type thirdPartyPkg struct {
	BaseUrl string
}

func NewThirdPartyPkg(baseUrl string) ThirdPartyPkg {
	return &thirdPartyPkg{BaseUrl: baseUrl}
}

func (t *thirdPartyPkg) CreateTransaction(accountId, reference string, amount float64) (*ThirdPartyPackageResponse, error) {
	reqBody := struct {
		AccountId string  `json:"account_id"`
		Reference string  `json:"reference"`
		Amount    float64 `json:"amount"`
	}{
		AccountId: accountId,
		Reference: reference,
		Amount:    amount,
	}
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/third-party/payments", t.BaseUrl), bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("failed to create payment: status code %d", resp.StatusCode)
	}

	var createdTransaction ThirdPartyPackageResponse
	if err := json.NewDecoder(resp.Body).Decode(&createdTransaction); err != nil {
		return nil, err
	}

	return &createdTransaction, nil
}

func (t *thirdPartyPkg) GetTransaction(reference string) (*ThirdPartyPackageResponse, error) {
	fullUrl := fmt.Sprintf("%s/third-party/payments/%s", t.BaseUrl, reference)

	req, err := http.NewRequest("GET", fullUrl, nil)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get payment: status code %d", resp.StatusCode)
	}

	var transaction ThirdPartyPackageResponse
	if err := json.NewDecoder(resp.Body).Decode(&transaction); err != nil {
		return nil, err
	}

	return &transaction, nil
}
