package telegram_api

import (
	"net/http"
	"net/url"
)

//go:generate mockery --name HTTPClientPost --output ./mocks
type HTTPClientPost interface {
	PostForm(url string, data url.Values) (*http.Response, error)
}

type WebHookReqBody struct {
	Message Message `json:"message"`
}

type Message struct {
	Text     string   `json:"text"`
	Chat     Chat     `json:"chat"`
	From     From     `json:"from"`
	Location Location `json:"location"`
}
type From struct {
	ID           int    `json:"id"`
	IsBot        bool   `json:"is_bot"`
	FirstName    string `json:"first_name"`
	Username     string `json:"username"`
	LanguageCode string `json:"language_code"`
}

type Chat struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name"`
	Username  string `json:"username"`
	Type      string `json:"type"`
}

type Config struct {
	Port       int    `env:"PORT" envDefault:"3000"`
	Token      string `env:"TOKEN"`
	API        string `env:"API"`
	Caller     int    `env:"CALLER"`
	TimeField  string `env:"TIMEFIELD"`
	TimeFormat string `env:"TIMEFORMAT"`
}

type KeyboardButton struct {
	Text     string `json:"text"`
	Location bool   `json:"request_location"`
}

type ReplyKeyboardMarkup struct {
	Keyboard [][]KeyboardButton `json:"keyboard"`
}

// Location struct for telegram body
type Location struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}
