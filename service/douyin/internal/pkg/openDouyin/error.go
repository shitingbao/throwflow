package openDouyin

import (
	"fmt"
)

type OpenDouyinError struct {
	Code        uint64
	Url         string
	Description string
	Response    string
}

func (ode *OpenDouyinError) Error() string {
	return fmt.Sprintf("[OpenDouyinError] Code=%d, Url=%s, Description=%s, Response=%s", ode.Code, ode.Url, ode.Description, ode.Response)
}

func NewOpenDouyinError(code uint64, url, description, response string) error {
	return &OpenDouyinError{
		Code:        code,
		Url:         url,
		Description: description,
		Response:    response,
	}
}

func (ode *OpenDouyinError) GetCode() uint64 {
	return ode.Code
}

func (ode *OpenDouyinError) GetUrl() string {
	return ode.Url
}

func (ode *OpenDouyinError) GetDescription() string {
	return ode.Description
}

func (ode *OpenDouyinError) GetResponse() string {
	return ode.Response
}
