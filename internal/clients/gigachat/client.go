package gigachat

import (
	"fmt"
	"io"
	"laboratory-internet-ai-test/config"
	"laboratory-internet-ai-test/internal/pkg/utils"
	"log/slog"
	"net/url"
	"strings"
	"sync"

	"context"
	"crypto/tls"
	"net/http"
	"time"
)

type Client struct {
	scope          string
	basicAuthToken string
	client         *http.Client
	accessToken    *accessToken
	mutex          *sync.Mutex

	logger *slog.Logger

	prompts map[string]config.Prompt
}

type accessToken struct {
	Token     string `json:"access_token"`
	ExpiresAt int64  `json:"expires_at"`
}

func New(config config.GigachatAiConfig, prompts map[string]config.Prompt, logger *slog.Logger) *Client {
	httpClient := Client{basicAuthToken: config.AuthKey, prompts: prompts, scope: config.Scope,
		mutex:  &sync.Mutex{},
		logger: logger,
		client: &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
			Timeout: 15 * time.Second,
		}}

	err := httpClient.getAccessToken(context.Background())
	if err != nil {
		logger.Error("Gigachat client init", "error", err)
	}

	return &httpClient
}

func (c *Client) ensureValidToken(ctx context.Context) error {

	c.mutex.Lock()
	defer c.mutex.Unlock()

	if c.accessToken.Token != "" && time.Now().Before(time.UnixMilli(c.accessToken.ExpiresAt)) {
		return nil
	}

	return c.getAccessToken(ctx)
}

func (c *Client) getAccessToken(ctx context.Context) error {

	data := url.Values{}
	data.Set("scope", c.scope)

	req, _ := http.NewRequestWithContext(ctx, "POST", "https://ngw.devices.sberbank.ru:9443/api/v2/oauth", strings.NewReader(data.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Basic "+c.basicAuthToken)
	req.Header.Set("RqUID", utils.GenerateUUID())

	resp, err := c.client.Do(req)
	if err != nil {
		c.logger.Error("getAccessToken gigachat", "error send request", err)
		return ErrGetToken
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK, http.StatusCreated:
	default:
		c.logger.Error("getAccessToken gigachat", "error", ErrRequestUnexpectedStatus, "status", resp.StatusCode)
		return ErrRequestUnexpectedStatus
	}

	bodyResp, err := io.ReadAll(resp.Body)
	if err != nil {
		c.logger.Error("getAccessToken gigachat", "error read body response", err)
		return ErrGetToken
	}

	response, err := utils.ParseJson[accessToken](bodyResp)
	if err != nil {
		c.logger.Error("getAccessToken gigachat", "error parse body response", err)
		return ErrGetToken
	}

	c.accessToken = response

	return nil
}

func (c *Client) generateBody(request *Request, command Command) (*RequestMessage, error) {

	comm := string(command)
	prompt, ok := c.prompts[comm]
	if !ok {
		return nil, fmt.Errorf("не найдена команда")
	}

	messageReq := RequestMessage{
		Model:             "GigaChat-2",
		N:                 1,
		Temperature:       0,
		Stream:            false,
		MaxTokens:         2024,
		RepetitionPenalty: 1,
		UpdateInterval:    0,
	}

	messages := make([]Message, 0, len(request.Messages)+1)
	messages = append(messages, Message{
		Role:    "system",
		Content: prompt.System,
	})

	for _, mess := range request.Messages {
		messages = append(messages, Message{
			Role:    "user",
			Content: mess,
		})
	}

	messageReq.Messages = messages

	if prompt.Function != nil {
		messageReq.FunctionGigachat = &FunctionGigachat{
			FunctionCall:    functionCall{Name: prompt.Name},
			Functions:       make([]Functions, 0, len(prompt.Function)),
			FewShotExamples: prompt.FewShotExamples,
		}
		for _, function := range prompt.Function {
			messageReq.FunctionGigachat.Functions = append(messageReq.FunctionGigachat.Functions, Functions{
				Name:        function.Name,
				Description: function.Description,
				Parameters: parameters{
					Type:       function.Parameters.Type,
					Properties: function.Parameters.Properties,
				},
				Required: function.Required,
			})
		}

	} else {
		messageReq.FunctionGigachat = &FunctionGigachat{
			FunctionCall:    functionCall{Name: "auto"},
			Functions:       nil,
			FewShotExamples: nil,
		}
	}

	return &messageReq, nil
}
