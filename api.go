package telegram_api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"sync"

	"github.com/phuslu/log"
)

var (
	lock      = &sync.Mutex{}
	singleApi *API
)

type API struct {
	client        HTTPClientPost
	url           string
	command       map[string]func(body *WebHookReqBody, chatID int) (url.Values, error)
	userInputFunc func(body *WebHookReqBody, chatID int) (url.Values, error)
}

func GetAPI(cfg *Config) *API {
	if singleApi == nil {
		lock.Lock()
		defer lock.Unlock()
		if singleApi == nil {
			singleApi = &API{
				client:  http.DefaultClient,
				url:     fmt.Sprintf("https://api.telegram.org/bot%v/sendMessage", cfg.Token),
				command: make(map[string]func(body *WebHookReqBody, chatID int) (url.Values, error)),
			}
			log.Info().Msg("API created")
		}
	}

	return singleApi
}

// Send response to a user
func (api *API) SendResponse(chatID int, val url.Values) error {
	response, err := api.client.PostForm(api.url, val)
	if err != nil {
		return fmt.Errorf("sending response failed. ChatID:%v, Text:%v.Error:%v", chatID, val.Get("text"), err)
	}

	if response == nil {
		return fmt.Errorf("response for ChatID:%v is nil. Input text:%v", chatID, val.Get("text"))
	}

	defer response.Body.Close()

	if response.StatusCode >= 400 && response.StatusCode < 500 {
		responseBody, readErr := io.ReadAll(response.Body)
		log.Warn().Msgf("Unable to send response for ChatID:%v. Text:%v. Response body:%v. Response error: %v", chatID, val.Get("text"), string(responseBody), readErr)

	}

	if response.StatusCode >= 500 {
		return fmt.Errorf("%v internal server error. ChatID:%v, Text:%s", response.StatusCode, chatID, val.Get("text"))
	}

	log.Debug().Msgf("Response to ChatID:%v sent. Message:%v. Response: %v", chatID, val.Get("text"), response)

	return nil
}

// TelegramHandler handles telegram request
func (api *API) TelegramHandler(_ http.ResponseWriter, r *http.Request) {
	var (
		body    *WebHookReqBody
		data    url.Values
		dataErr error
	)

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		log.Error().Err(err).Msgf("error occurred decoding message body: %v. Error:%v", r.Body, err)
	}

	chatID := strconv.Itoa(body.Message.Chat.ID)
	textToSend := url.Values{"chat_id": {chatID}, "text": {fmt.Sprintf("unable to find a response for %s", body.Message.Text)}}

	function, found := api.command[body.Message.Text]
	if !found {
		data, dataErr = api.userInputFunc(body, body.Message.Chat.ID)
		if dataErr != nil {
			api.SendResponse(body.Message.Chat.ID, textToSend)
			log.Error().Err(fmt.Errorf("unable to find response function for %s", dataErr))
			return
		}
		api.SendResponse(body.Message.Chat.ID, data)
		return
	}

	data, dataErr = function(body, body.Message.Chat.ID)
	if dataErr != nil {
		log.Error().Err(fmt.Errorf("an error occurred while attempting to retrieve an answer for the user. See error: %w", dataErr))
	}

	err := api.SendResponse(body.Message.Chat.ID, data)
	if err != nil {
		log.Error().Err(fmt.Errorf("SendResponse error: %w", err))
		api.SendResponse(body.Message.Chat.ID, url.Values{"chat_id": {chatID}, "text": {fmt.Sprintf("SendResponse error: %v", err)}})
		return
	}
}

func (api *API) RegisterCommand(command string, callback func(body *WebHookReqBody, chatID int) (url.Values, error)) {
	api.command[command] = callback
}
func (api *API) RegisterInput(callback func(body *WebHookReqBody, chatID int) (url.Values, error)) {
	api.userInputFunc = callback
}
