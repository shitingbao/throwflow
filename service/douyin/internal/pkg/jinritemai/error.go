package jinritemai

import (
	"fmt"
)

type JinritemaiError struct {
	Code        uint64
	Url         string
	Description string
	Response    string
}

func (je *JinritemaiError) Error() string {
	return fmt.Sprintf("[JinritemaiError] Code=%d, Url=%s, Description=%s, Response=%s", je.Code, je.Url, je.Description, je.Response)
}

func NewJinritemaiError(code uint64, url, description, response string) error {
	return &JinritemaiError{
		Code:        code,
		Url:         url,
		Description: description,
		Response:    response,
	}
}

func (je *JinritemaiError) GetCode() uint64 {
	return je.Code
}

func (je *JinritemaiError) GetUrl() string {
	return je.Url
}

func (je *JinritemaiError) GetDescription() string {
	return je.Description
}

func (je *JinritemaiError) GetResponse() string {
	return je.Response
}
