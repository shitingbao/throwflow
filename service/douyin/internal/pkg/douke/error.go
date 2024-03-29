package douke

type DoukeError struct {
	Code        uint64
	Url         string
	Description string
	Response    string
}

func (de *DoukeError) Error() string {
	return de.Description
}

func NewDoukeError(code uint64, url, description, response string) error {
	return &DoukeError{
		Code:        code,
		Url:         url,
		Description: description,
		Response:    response,
	}
}

func (de *DoukeError) GetCode() uint64 {
	return de.Code
}

func (de *DoukeError) GetUrl() string {
	return de.Url
}

func (de *DoukeError) GetDescription() string {
	return de.Description
}

func (de *DoukeError) GetResponse() string {
	return de.Response
}
