package openDouyin

type CommonExtraResponse struct {
	ErrorCode      uint64 `json:"error_code"`
	SubErrorCode   uint64 `json:"sub_error_code"`
	Description    string `json:"description"`
	SubDescription string `json:"sub_description"`
	Logid          string `json:"logid"`
	Now            uint64 `json:"now"`
}

type CommonResponse struct {
	Extra CommonExtraResponse `json:"extra"`
}
