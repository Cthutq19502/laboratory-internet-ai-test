package gigachat

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	domaincontact "laboratory-internet-ai-test/internal/domain/contact"
	"laboratory-internet-ai-test/internal/pkg/utils"
	"net/http"
)

func (c *Client) GetTonal(ctx context.Context, request *Request) (domaincontact.Tonal, error) {

	requestBody, err := c.generateBody(request, GetTonalCommand)
	if err != nil {
		c.logger.Error("GetTonal gigachat", "error generate gigachat body request", err)
		return "", err
	}

	response, err := c.sendMessageOnce(ctx, requestBody)
	if err != nil {
		c.logger.Error("GetTonal gigachat", "error send gigachat request", err)
		return "", err
	}

	if response.Choices == nil || len(response.Choices) == 0 {
		c.logger.Error("GetTonal gigachat", "error", "choices is null")
		return "", fmt.Errorf("choices is null")
	}
	functionCall := response.Choices[0].MessageResponse.FunctionCall
	args := functionCall.Arguments
	if arg, exist := args["tonal"]; exist {
		tonal, ok := (arg).(string)
		if !ok {
			c.logger.Error("GetTonal gigachat", "error", "wrong arguments from gigachat", "args", args)
			return "", fmt.Errorf("unexpected arguments tonal response")
		}
		t := domaincontact.Tonal(tonal)
		switch t {
		case domaincontact.TonalNeutral, domaincontact.TonalNegative, domaincontact.TonalPositive:
			return t, nil
		default:
			c.logger.Error("GetTonal gigachat", "error", "unexpected tonal response", "response", tonal)
			return "", fmt.Errorf("unexpected tonal response")
		}
	}

	return domaincontact.TonalUnexpected, nil

}

func (c *Client) sendMessageOnce(ctx context.Context, message *RequestMessage) (*Response, error) {

	if err := c.ensureValidToken(ctx); err != nil {
		return nil, ErrTokenExpired
	}

	bodyReq, err := json.Marshal(message)
	if err != nil {
		c.logger.Error("SendMessageOnce gigachat", "error marshal body", err)
		return nil, ErrSendMessage
	}

	c.logger.InfoContext(ctx, "SendMessageOnce gigachat", "generate body:", bodyReq)

	req, _ := http.NewRequest("POST", "https://gigachat.devices.sberbank.ru/api/v1/chat/completions", bytes.NewBuffer(bodyReq))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.accessToken.Token)

	resp, err := c.client.Do(req)
	if err != nil {
		c.logger.Error("SendMessageOnce gigachat", "error request", err)
		if errors.Is(err, context.DeadlineExceeded) {
			return nil, ErrRequestTimeout
		}
		return nil, ErrSendMessage
	}

	defer resp.Body.Close()
	bodyResp, err := io.ReadAll(resp.Body)
	if err != nil {
		c.logger.Error("SendMessageOnce gigachat", "error read body response", err)
		return nil, ErrSendMessage
	}

	switch resp.StatusCode {
	case http.StatusCreated, http.StatusOK:
	default:
		c.logger.Error("SendMessageOnce gigachat", "error", ErrRequestUnexpectedStatus, "status", resp.StatusCode,
			"response", bodyResp)
		return nil, ErrRequestUnexpectedStatus
	}

	response, err := utils.ParseJson[Response](bodyResp)
	if err != nil {
		c.logger.Error("SendMessageOnce gigachat", "error parse body response", err, "body", bodyResp)
		return nil, ErrSendMessage
	}

	c.logger.InfoContext(ctx, "SendMessageOnce gigachat", "response body:", response)

	return response, nil
}
