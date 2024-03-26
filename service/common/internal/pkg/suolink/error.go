package suolink

import (
	"fmt"
)

type SuoLinkError struct {
	Url         string
	Description string
	Response    string
}

func (sle *SuoLinkError) Error() string {
	return fmt.Sprintf("[SuoLinkError] Url=%s, Description=%s, Response=%s", sle.Url, sle.Description, sle.Response)
}

func NewSuoLinkError(url, description, response string) error {
	return &SuoLinkError{
		Url:         url,
		Description: description,
		Response:    response,
	}
}

func (sle *SuoLinkError) GetUrl() string {
	return sle.Url
}

func (sle *SuoLinkError) GetDescription() string {
	return sle.Description
}

func (sle *SuoLinkError) GetResponse() string {
	return sle.Response
}
