package csj

type CsjError struct {
	Code        int64
	Url         string
	Description string
	Response    string
}

func (ce *CsjError) Error() string {
	return ce.Description
}

func NewCsjError(code int64, url, description, response string) error {
	return &CsjError{
		Code:        code,
		Url:         url,
		Description: description,
		Response:    response,
	}
}

func (ce *CsjError) GetCode() int64 {
	return ce.Code
}

func (ce *CsjError) GetUrl() string {
	return ce.Url
}

func (ce *CsjError) GetDescription() string {
	return ce.Description
}

func (ce *CsjError) GetResponse() string {
	return ce.Response
}
