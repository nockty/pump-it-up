package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/nockty/pump-it-up/internal/telegram"
)

type secrets struct {
	APIKey    string `json:"binanceApiKey"`
	SecretKey string `json:"binanceSecretKey"`
	APIID     string `json:"telegramApiId"`
	APIHash   string `json:"telegramApiHash"`
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
	// client := binance.NewClient(secrets.APIKey, secrets.SecretKey)
	// res, err := client.NewGetAccountService().Do(context.Background())
	// if err != nil {
	// 	println(err)
	// 	return
	// }
	// for _, b := range res.Balances {
	// 	if b.Asset != "BTC" {
	// 		continue
	// 	}
	// 	fmt.Printf("Balance: %s\n", b.Free)
	// }
	tp := telegram.NewTelegramPoller(secrets.APIID, secrets.APIHash)
	tp.GetInteractiveAuthorization()
	go tp.Run()
	for coin := range tp.BuyChan {
		fmt.Printf("Buy %s now!", coin)
	}
}
