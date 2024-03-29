package oceanengine

import (
	"fmt"
)

type OceanengineError struct {
	Code        uint32
	Url         string
	Message     string
	Description string
	Response    string
	RequestId   string
}

func (oe *OceanengineError) Error() string {
	if oe.RequestId == "" {
		return fmt.Sprintf("[OceanengineError] Code=%d, Url=%s, Message=%s, Description=%s, Response=%s", oe.Code, oe.Url, oe.Message, oe.Description, oe.Response)
	}

	return fmt.Sprintf("[OceanengineError] Code=%d, Url=%s, Message=%s, Description=%s, RequestId=%s, Response=%s", oe.Code, oe.Url, oe.Message, oe.Description, oe.RequestId, oe.Response)
}

func NewOceanengineError(code uint32, url, message, description, requestId, response string) error {
	return &OceanengineError{
		Code:        code,
		Url:         url,
		Message:     message,
		Description: description,
		Response:    response,
		RequestId:   requestId,
	}
}

func (oe *OceanengineError) GetCode() uint32 {
	return oe.Code
}

func (oe *OceanengineError) GetUrl() string {
	return oe.Url
}

func (oe *OceanengineError) GetMessage() string {
	return oe.Message
}

func (oe *OceanengineError) GetDescription() string {
	return oe.Description
}

func (oe *OceanengineError) GetResponse() string {
	return oe.Response
}

func (oe *OceanengineError) GetRequestId() string {
	return oe.RequestId
}
