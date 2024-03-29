package jinritemai

import (
	"crypto/sha1"
	"encoding/hex"
)

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

func VerifyMessageSign(appSecret, paramJson, sign string) bool {
	hash := sha1.New()
	hash.Write([]byte(appSecret + paramJson))

	return sign == hex.EncodeToString(hash.Sum(nil))
}

type MessageResponse struct {
	Event      string `json:"event"`
	ClientKey  string `json:"client_key"`
	FromUserId string `json:"from_user_id"`
	Content    string `json:"content"`
	LogId      string `json:"log_id"`
}
