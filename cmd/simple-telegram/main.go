package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Arman92/go-tdlib"
)

func main() {
	tdlib.SetLogVerbosityLevel(1)
	tdlib.SetFilePath("./errors.txt")

	// Create new instance of client
	client := tdlib.NewClient(tdlib.Config{
		APIID:               "187786",
		APIHash:             "e782045df67ba48e441ccb105da8fc85",
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
	})

	// Handle Ctrl+C, gracefully exit and shutdown tdlib
	ch := make(chan os.Signal, 2)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-ch
		client.DestroyInstance()
		os.Exit(1)
	}()

	for {
		currentState, _ := client.Authorize()
		if currentState.GetAuthorizationStateEnum() == tdlib.AuthorizationStateWaitPhoneNumberType {
			fmt.Print("Enter phone: ")
			var number string
			fmt.Scanln(&number)
			_, err := client.SendPhoneNumber(number)
			if err != nil {
				fmt.Printf("Error sending phone number: %v", err)
			}
		} else if currentState.GetAuthorizationStateEnum() == tdlib.AuthorizationStateWaitCodeType {
			fmt.Print("Enter code: ")
			var code string
			fmt.Scanln(&code)
			_, err := client.SendAuthCode(code)
			if err != nil {
				fmt.Printf("Error sending auth code : %v", err)
			}
		} else if currentState.GetAuthorizationStateEnum() == tdlib.AuthorizationStateWaitPasswordType {
			fmt.Print("Enter Password: ")
			var password string
			fmt.Scanln(&password)
			_, err := client.SendAuthPassword(password)
			if err != nil {
				fmt.Printf("Error sending auth password: %v", err)
			}
		} else if currentState.GetAuthorizationStateEnum() == tdlib.AuthorizationStateReadyType {
			fmt.Println("Authorization Ready! Let's rock")
			break
		}
	}

	go func() {
		// Create a filter function which will be used to filter out unwanted tdlib messages
		eventFilter := func(msg *tdlib.TdMessage) bool {
			return true
			// 	updateMsg := (*msg).(*tdlib.UpdateNewMessage)
			// 	// For example, we want incomming messages from user with below id:
			// 	if updateMsg.Message.SenderUserID == 41507975 {
			// 		retufalsern true
			// 	}
			// 	return false
		}

		// Here we can add a receiver to retreive any message type we want
		// We like to get UpdateNewMessage events and with a specific FilterFunc
		receiver := client.AddEventReceiver(&tdlib.UpdateNewMessage{}, eventFilter, 5)
		for newMsg := range receiver.Chan {
			fmt.Println(newMsg)
			updateMsg := (newMsg).(*tdlib.UpdateNewMessage)
			// We assume the message content is simple text: (should be more sophisticated for general use)
			msgText := updateMsg.Message.Content.(*tdlib.MessageText)
			fmt.Println("MsgText:  ", msgText.Text)
			fmt.Print("\n\n")
		}

	}()

	for {
		time.Sleep(10)
	}

	// rawUpdates gets all updates comming from tdlib
	// rawUpdates := client.GetRawUpdatesChannel(100)
	// for update := range rawUpdates {
	// 	// Show all updates
	// 	fmt.Println(update.Data)
	// 	fmt.Print("\n\n")
	// }

}
