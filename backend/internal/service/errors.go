package service

import "errors"

var (
	ErrLLMProviderFailed = errors.New("llm provider failed")
	ErrInvalidLLMOutput  = errors.New("invalid llm output")
)
