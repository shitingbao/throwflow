package kuaidi

import (
	"fmt"
)

type KuaidiError struct {
	Url         string
	Description string
	Response    string
}

func (ke *KuaidiError) Error() string {
	return fmt.Sprintf("[KuaidiError] Url=%s, Description=%s, Response=%s", ke.Url, ke.Description, ke.Response)
}

func NewKuaidiError(url, description, response string) error {
	return &KuaidiError{
		Url:         url,
		Description: description,
		Response:    response,
	}
}

func (ke *KuaidiError) GetUrl() string {
	return ke.Url
}

func (ke *KuaidiError) GetDescription() string {
	return ke.Description
}

func (ke *KuaidiError) GetResponse() string {
	return ke.Response
}
