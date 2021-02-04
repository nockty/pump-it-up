package main

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/nockty/pump-it-up/internal/binance"
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
	args := os.Args
	if len(args) != 2 {
		println("Usage: simple-binance [BTCAmount]")
		os.Exit(64)
	}
	BTCAmount := args[1]
	secrets, err := getSecrets(".secrets.json")
	if err != nil {
		println(err)
		return
	}
	bt := binance.NewTrader(secrets.APIKey, secrets.SecretKey, BTCAmount)
	tp := telegram.NewPoller(secrets.APIID, secrets.APIHash)
	tp.GetInteractiveAuthorization()
	go tp.Run()
	coin := <-tp.BuyChan
	bt.Trade(coin)
}
