# Telegram API Package

This Go package provides a convenient way to interact with the Telegram Bot API. It includes functionalities for handling incoming Telegram webhook requests, sending responses to users, and registering custom commands.

## Installation

To use this package, you can install it using the following `go get` command:

```bash
go get "github.com/c1kzy/Telegram-API"
```

## Usage

``` go
package main

import (
	"fmt"
	"github.com/c1kzy/Telegram-API"
)

func main() {
    // Initialize the configuration
    cfg := &telegram_api.Config{
        Token: "YOUR_TELEGRAM_BOT_TOKEN",
    }

    // Get the API instance
    api := telegram_api.GetAPI(cfg)

    // Register custom commands
    api.RegisterCommand("/start", func(body *telegram_api.WebHookReqBody, chatID int) (url.Values, error) {
        // Your custom command logic here
        return url.Values{"chat_id": {fmt.Sprintf("%d", chatID)}, "text": {"Hello, welcome to the bot!"}}, nil
    })

    // Register custom user input handler
    api.RegisterInput(func(body *telegram_api.WebHookReqBody, chatID int) (url.Values, error) {
        // Your custom user input logic here
        return url.Values{"chat_id": {fmt.Sprintf("%d", chatID)}, "text": {"Sorry, I couldn't understand your request."}}, nil
    })

    // Start the HTTP server to handle incoming Telegram webhook requests
    http.HandleFunc("/webhook", api.TelegramHandler)
    if err := http.ListenAndServe(":8080", nil); err != nil {
        fmt.Println("Error starting the server:", err)
    }
}
