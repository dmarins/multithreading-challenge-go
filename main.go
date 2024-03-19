package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

type PayloadViaCepApi struct {
	Cep         string `json:"cep"`
	Logradouro  string `json:"logradouro"`
	Complemento string `json:"complemento"`
	Bairro      string `json:"bairro"`
	Localidade  string `json:"localidade"`
	Uf          string `json:"uf"`
	Ibge        string `json:"ibge"`
	Gia         string `json:"gia"`
	Ddd         string `json:"ddd"`
	Siafi       string `json:"siafi"`
}

type PayloadBrasilApi struct {
	Cep          string `json:"cep"`
	Street       string `json:"street"`
	Neighborhood string `json:"neighborhood"`
	City         string `json:"city"`
	State        string `json:"state"`
	Service      string `json:"service"`
}

func main() {
	ctx := context.Background()

	c1 := make(chan PayloadViaCepApi)
	c2 := make(chan PayloadBrasilApi)

	go func() {
		viaCep, err := getCepWithViaCepApi(ctx)
		if err != nil {
			log.Println(err)
			return
		}
		c1 <- *viaCep
	}()

	go func() {
		brasilApi, err := getCepWithBrasilApi(ctx)
		if err != nil {
			log.Println(err)
			return
		}
		c2 <- *brasilApi
	}()

	select {
	case cep := <-c1:
		fmt.Printf("Via CEP API: %s\n", cep)
	case cep := <-c2:
		fmt.Printf("Brasil API: %s\n", cep)
	case <-time.After(time.Second * 1):
		fmt.Printf("timeout.\n")

	}
}

func getCepWithViaCepApi(ctx context.Context) (*PayloadViaCepApi, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)

	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", "https://viacep.com.br/ws/24452050/json/", nil)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	var rawJson map[string]interface{}
	json.Unmarshal(body, &rawJson)

	strJson, _ := json.Marshal(rawJson)

	var payloadViaCepApi PayloadViaCepApi
	err = json.Unmarshal(strJson, &payloadViaCepApi)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &payloadViaCepApi, nil
}

func getCepWithBrasilApi(ctx context.Context) (*PayloadBrasilApi, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)

	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", "https://brasilapi.com.br/api/cep/v1/24452050", nil)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	var rawJson map[string]interface{}
	json.Unmarshal(body, &rawJson)

	strJson, _ := json.Marshal(rawJson)

	var payloadBrasilApi PayloadBrasilApi
	err = json.Unmarshal(strJson, &payloadBrasilApi)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &payloadBrasilApi, nil
}
