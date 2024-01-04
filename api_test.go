package telegram_api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"testing"
	"time"

	"github.com/phuslu/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockHTTPClient struct {
	mock.Mock
}

func (c *MockHTTPClient) PostForm(url string, data url.Values) (*http.Response, error) {
	args := c.Called(url, data)
	return args.Get(0).(*http.Response), args.Error(1)
}

func marshJSON(t *testing.T, text string, chatID int) string {
	reqBody := &WebHookReqBody{
		Message: Message{
			Text: text,
			Chat: Chat{
				ID: chatID,
			},
		},
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		t.Fatal(err)
	}
	return string(body)
}

func TestHandlerTelegram(t *testing.T) {
	var cfg Config

	log.DefaultLogger = log.Logger{
		Level:      log.DebugLevel,
		Caller:     cfg.Caller,
		TimeField:  cfg.TimeField,
		TimeFormat: time.RFC850,
		Writer:     &log.ConsoleWriter{},
	}

	recorder := httptest.NewRecorder()

	api := GetAPI(&cfg)
	log.Debug().Msgf("Config loaded: %v", cfg)
	api.RegisterCommand("/start", func(body *WebHookReqBody, chatID int) (url.Values, error) {
		return url.Values{
			"chat_id": {strconv.Itoa(chatID)},
			"text":    {"test start case"},
		}, nil
	})
	api.RegisterInput(func(body *WebHookReqBody, chatID int) (url.Values, error) {
		return url.Values{
			"chat_id": {strconv.Itoa(chatID)},
			"text":    {fmt.Sprintf("This is your input: %v\n", body.Message.Text)},
		}, nil
	})

	tests := []struct {
		name    string
		request *http.Request
	}{
		{
			name:    "/start",
			request: httptest.NewRequest(http.MethodPost, "/telegram", bytes.NewBuffer([]byte(marshJSON(t, "/start", 123)))),
		},
		{
			name:    "test input",
			request: httptest.NewRequest(http.MethodPost, "/telegram", bytes.NewBuffer([]byte(marshJSON(t, "test input", 123)))),
		},

		{
			name:    "invalid command",
			request: httptest.NewRequest(http.MethodPost, "/telegram", bytes.NewBuffer([]byte(marshJSON(t, "/invalid", 123)))),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			api.TelegramHandler(recorder, tt.request)
			assert.Equal(t, http.StatusOK, recorder.Code)
		})
	}
}
