package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jasonmoo/lambda_proc"
	"io/ioutil"
	"net/http"
)

type Event struct {
	Ticker string `json:"ticker"`
}
type Response struct {
	Price  float64
	Volume string
	Name   string
	Symbol string
}

type YahooResult struct {
	List struct {
		Resources []struct {
			Resource struct {
				Fields struct {
					Price  float64 `json:"price,string"`
					Volume string  `json:"volume"`
					Name   string  `json:"name"`
					Symbol string  `json:"symbol"`
				} `json:"fields"`
			} `json:"resource"`
		} `json:"resources"`
	} `json:"list"`
}

func main() {
	lambda_proc.Run(HandleEvent)
}

func HandleEvent(context *lambda_proc.Context, eventJSON json.RawMessage) (interface{}, error) {
	event := &Event{}
	data := &YahooResult{}
	response := Response{}
	if err := json.Unmarshal(eventJSON, event); err != nil {
		return nil, err
	}
	if event.Ticker == "" {
		return nil, errors.New("ticker required")
	}

	endpoint := fmt.Sprintf(`http://finance.yahoo.com/webservice/v1/symbols/%s/quote?format=json`, event.Ticker)
	resp, err := http.Get(endpoint)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err := json.Unmarshal(body, data); err != nil {
		return nil, err
	}
	if len(data.List.Resources) < 1 {
		return nil, errors.New("ticker not on yahoo")
	}
	response.Price = data.List.Resources[0].Resource.Fields.Price
	response.Volume = data.List.Resources[0].Resource.Fields.Volume
	response.Name = data.List.Resources[0].Resource.Fields.Name
	response.Symbol = data.List.Resources[0].Resource.Fields.Symbol
	return response, nil
}
