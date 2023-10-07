package exchangerates

import (
	"encoding/json"
	"fmt"
	"github.com/groshi-project/groshi/internal/loggers"
	"io"
	"net/http"
)

type apiClient struct {
	accessKey  string
	httpClient http.Client
}

func newAPIClient() *apiClient {
	return &apiClient{
		httpClient: http.Client{},
	}
}

func (client *apiClient) Init(accessKey string) {
	client.accessKey = accessKey
}

func (client *apiClient) sendRequest(url string) (map[string]interface{}, error) {
	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	response, err := client.httpClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := response.Body.Close(); err != nil {
			loggers.Error.Printf("failed to close response body: %v", err)
		}
	}()
	responseBody, err := io.ReadAll(response.Body)

	var responseJSON map[string]interface{}
	if err := json.Unmarshal(responseBody, &responseJSON); err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK || responseJSON["success"] != true {
		return nil, fmt.Errorf(
			"non-successful response from the server: %v", string(responseBody),
		)
	}

	return responseJSON, nil
}

// GetRates returns rates for provided baseCurrency.
func (client *apiClient) GetRates(baseCurrency string) (map[string]interface{}, error) {
	url := fmt.Sprintf(
		"http://api.exchangeratesapi.io/v1/latest?access_key=%v&base=%v",
		client.accessKey,
		baseCurrency,
	)
	response, err := client.sendRequest(url)
	if err != nil {
		return nil, err
	}
	return response["rates"].(map[string]interface{}), nil
}

// GetCurrencies returns supported currency codes (ISO-4217).
func (client *apiClient) GetCurrencies() ([]string, error) {
	url := fmt.Sprintf("http://api.exchangeratesapi.io/v1/symbols?access_key=%v", client.accessKey)
	response, err := client.sendRequest(url)
	if err != nil {
		return nil, err
	}
	symbols, _ := response["symbols"].(map[string]interface{})

	var keys []string
	for key, value := range symbols {
		fmt.Printf("%v:%v", key, value)
		keys = append(keys, key)
	}
	return keys, err
}

var Client = newAPIClient()
