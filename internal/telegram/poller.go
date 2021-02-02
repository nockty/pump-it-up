package telegram

import (
	"fmt"
	"regexp"

	"github.com/Arman92/go-tdlib"
)

// const bigPumpSignalID int64 = -1001257721998
const bigPumpSignalID int64 = 260737464
const pumpCoinRegex = `^.*\$([A-Z]+).*$`

type TelegramPoller struct {
	BuyChan chan string
	client  *tdlib.Client
	regex   *regexp.Regexp
}

func NewTelegramPoller(APIID, APIHash string) *TelegramPoller {
	tdlib.SetLogVerbosityLevel(1)
	tdlib.SetFilePath("./errors.txt")

	regex, _ := regexp.Compile(pumpCoinRegex)

	return &TelegramPoller{
		BuyChan: make(chan string, 10),
		client: tdlib.NewClient(tdlib.Config{
			APIID:               APIID,
			APIHash:             APIHash,
			SystemLanguageCode:  "en",
			DeviceModel:         "Server",
			SystemVersion:       "1.0.0",
			ApplicationVersion:  "1.0.0",
			UseMessageDatabase:  true,
			UseFileDatabase:     true,
			UseChatInfoDatabase: true,
			UseTestDataCenter:   false,
			DatabaseDirectory:   "./tdlib-db",
			FileDirectory:       "./tdlib-files",
			IgnoreFileNames:     false,
		}),
		regex: regex,
	}
}

func (tp *TelegramPoller) GetInteractiveAuthorization() {
	for {
		currentState, _ := tp.client.Authorize()
		if currentState.GetAuthorizationStateEnum() == tdlib.AuthorizationStateWaitPhoneNumberType {
			println("Enter phone: ")
			var number string
			fmt.Scanln(&number)
			_, err := tp.client.SendPhoneNumber(number)
			if err != nil {
				fmt.Printf("Error sending phone number: %v", err)
			}
		} else if currentState.GetAuthorizationStateEnum() == tdlib.AuthorizationStateWaitCodeType {
			print("Enter code: ")
			var code string
			fmt.Scanln(&code)
			_, err := tp.client.SendAuthCode(code)
			if err != nil {
				fmt.Printf("Error sending auth code : %v", err)
			}
		} else if currentState.GetAuthorizationStateEnum() == tdlib.AuthorizationStateWaitPasswordType {
			print("Enter Password: ")
			var password string
			fmt.Scanln(&password)
			_, err := tp.client.SendAuthPassword(password)
			if err != nil {
				fmt.Printf("Error sending auth password: %v", err)
			}
		} else if currentState.GetAuthorizationStateEnum() == tdlib.AuthorizationStateReadyType {
			println("Authorization Ready! Let's rock")
			break
		}
	}
}

func (tp *TelegramPoller) Run() {
	eventFilter := func(msg *tdlib.TdMessage) bool {
		updateMsg := (*msg).(*tdlib.UpdateNewMessage)
		if updateMsg.Message.ChatID == bigPumpSignalID {
			return true
		}
		return false
	}

	receiver := tp.client.AddEventReceiver(&tdlib.UpdateNewMessage{}, eventFilter, 5)
	for newMsg := range receiver.Chan {
		fmt.Println(newMsg)
		updateMsg := (newMsg).(*tdlib.UpdateNewMessage)
		// We assume the message content is simple text: (should be more sophisticated for general use)
		msgText := updateMsg.Message.Content.(*tdlib.MessageText).Text.Text
		submatch := tp.regex.FindStringSubmatch(msgText)
		if submatch != nil {
			tp.BuyChan <- submatch[1]
			close(tp.BuyChan)
		}
	}
}
