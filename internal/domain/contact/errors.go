package contact

import "errors"

var ErrGetFileUUID = errors.New("error get file uuid")
var ErrEmptyResponse = errors.New("empty response from gigachat")
var ErrBlacklistResponse = errors.New("blacklist response from gigachat")
var ErrGenerateMessage = errors.New("error generate gigachat message")
var ErrCommandNotFound = errors.New("command not found")

//----------------------------------------------------------------------

var ErrInvalidInput = errors.New("invalid input")
