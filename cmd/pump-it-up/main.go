package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/adshao/go-binance/v2"
)

type secrets struct {
	APIKey    string `json:"apiKey"`
	SecretKey string `json:"secretKey"`
}

func getSecrets(fileName string) (*secrets, error) {
	secretsFile, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer secretsFile.Close()
	byteValue, err := ioutil.ReadAll(secretsFile)
	if err != nil {
		return nil, err
	}
	var secrets secrets
	json.Unmarshal(byteValue, &secrets)
	return &secrets, nil
}

func main() {
	secrets, err := getSecrets(".secrets.json")
	if err != nil {
		println(err)
		return
	}
	client := binance.NewClient(secrets.APIKey, secrets.SecretKey)
	res, err := client.NewGetAccountService().Do(context.Background())
	if err != nil {
		println(err)
		return
	}
	for _, b := range res.Balances {
		if b.Asset != "BTC" {
			continue
		}
		fmt.Printf("Balance: %s\n", b.Free)
	}
}
