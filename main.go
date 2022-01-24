package main

import (
	"fmt"
	"log"
	"os"
	"time"

	monobank "github.com/RomanMelnyk113/monobank-sdk"
	tele "gopkg.in/telebot.v3"
)

func main() {
	users := make(map[int64]*monobank.UserInfo)
	token := os.Getenv("TOKEN")
	pref := tele.Settings{
		Token:  token,
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	b, err := tele.NewBot(pref)
	if err != nil {
		log.Fatal(err)
		return
	}

	b.Handle("/hello", func(c tele.Context) error {
		msg := "Hello! Check available commands below: \n\n"
		start := "/start <token> - provide me your Monobank API key so I can load your profile details.\n\n"
		accounts := "/accounts - display list of available accounts.\n\n"
		return c.Send(msg + start + accounts)
	})

	// Accept monobank token to start conversation
	// Example: /start <token>
	b.Handle("/start", func(c tele.Context) error {
		token := c.Message().Payload
		client := monobank.NewClient(token)
		user, err := client.GetUserInfo()
		if err != nil {
			log.Fatal(err)
			return err
		}
		users[c.Sender().ID] = user
		msg := "Hello " + user.Name + ". Run /accounts to check available accounts"
		return c.Send(msg)
	})

	// Return list of available accounts
	b.Handle("/accounts", func(c tele.Context) error {
		var msg string = ""
		user, exists := users[c.Sender().ID]
		if !exists {
			return c.Send("Please call /start <monobank_token> first")
		}
		// Add support to convert CurrencyCode to CurrencySign
		for _, account := range user.Accounts {
			msg += fmt.Sprintf("AccountID: %s, Balance %d, CurrencyCode %d\n\n", account.AccountID, account.Balance, account.CurrencyCode)
		}
		return c.Send(msg)
	})

	b.Handle(tele.OnText, func(c tele.Context) error {
		// All the text messages that weren't
		// captured by existing handlers.

		var (
			user = c.Sender()
			text = c.Text()
		)

		// Use full-fledged bot's functions
		// only if you need a result:
		msg, err := b.Send(user, text)
		if err != nil {
			return err
		}

		// Instead, prefer a context short-hand:
		return c.Send(msg)
	})

	b.Start()
}
