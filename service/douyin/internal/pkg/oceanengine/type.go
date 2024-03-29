package oceanengine

type CommonResponse struct {
	Code      uint32 `json:"code"`
	Message   string `json:"message"`
	RequestId string `json:"request_id"`
}

type PageResponse struct {
	Page        uint32 `json:"page"`
	PageSize    uint32 `json:"page_size"`
	TotalNumber uint32 `json:"total_number"`
	TotalPage   uint32 `json:"total_page"`
}

type CursorPageResponse struct {
	Cursor  uint32 `json:"cursor"`
	HasMore bool   `json:"has_more"`
}

type TokenResponse struct {
	AccessToken           string `json:"access_token"`
	ExpiresIn             uint32 `json:"expires_in"`
	RefreshToken          string `json:"refresh_token"`
	RefreshTokenExpiresIn uint32 `json:"refresh_token_expires_in"`
}
