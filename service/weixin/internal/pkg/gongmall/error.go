package gongmall

type GongmallError struct {
	Code        string
	Url         string
	Description string
	Response    string
}

func (ge *GongmallError) Error() string {
	return ge.Description
}

func NewGongmallError(code, url, description, response string) error {
	return &GongmallError{
		Code:        code,
		Url:         url,
		Description: description,
		Response:    response,
	}
}

func (ge *GongmallError) GetCode() string {
	return ge.Code
}

func (ge *GongmallError) GetUrl() string {
	return ge.Url
}

func (ge *GongmallError) GetDescription() string {
	return ge.Description
}

func (ge *GongmallError) GetResponse() string {
	return ge.Response
}
