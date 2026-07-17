package gigachat

import "errors"

var ErrTokenExpired = errors.New("error generate token")
var ErrRequestTimeout = errors.New("request waiting timeout")
var ErrRequestUnexpectedStatus = errors.New("request unexpected status")
var ErrSendMessage = errors.New("send message request error")
var ErrGetToken = errors.New("get token request error")
