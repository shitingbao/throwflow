package douke

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"douyin/internal/pkg/tool"
	"encoding/hex"
	"encoding/json"
	"github.com/google/go-querystring/query"
	"strings"
)

type CommonRequest struct {
	Method      string `url:"method"`
	AppKey      string `url:"app_key"`
	AccessToken string `url:"access_token"`
	ParamJson   string `url:"param_json"`
	Timestamp   string `url:"timestamp"`
	V           string `url:"v"`
	Sign        string `url:"sign"`
	SignMethod  string `url:"sign_method"`
}

func (cr CommonRequest) String() string {
	v, _ := query.Values(cr)

	return v.Encode()
}

func Sign(appKey, appSecret, method, paramJson, timestamp, v string) string {
	paramPattern := "app_key" + appKey + "method" + method + "param_json" + paramJson + "timestamp" + timestamp + "v" + v

	signPattern := appSecret + paramPattern + appSecret

	return Hmac(signPattern, appSecret)
}

func VerifyMessageSign(appKey, appSecret, paramJson, sign string) bool {
	signPattern := appKey + paramJson + appSecret

	return sign == tool.GetMd5(signPattern)
}

func Hmac(s string, appSecret string) string {
	h := hmac.New(sha256.New, []byte(appSecret))
	_, _ = h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

func Marshal(o interface{}) string {
	raw, _ := json.Marshal(o)

	m := make(map[string]interface{})
	reader := bytes.NewReader(raw)
	decode := json.NewDecoder(reader)
	decode.UseNumber()
	_ = decode.Decode(&m)

	buffer := bytes.NewBufferString("")
	encoder := json.NewEncoder(buffer)
	encoder.SetEscapeHTML(false)
	_ = encoder.Encode(m)

	marshal := strings.TrimSpace(buffer.String())
	return marshal
}

type CommonResponse struct {
	Code    uint64 `json:"code"`
	Msg     string `json:"msg"`
	SubCode string `json:"sub_code"`
	SubMsg  string `json:"sub_msg"`
	LogId   string `json:"log_id"`
}

type MessageResponse struct {
	Tag   string `json:"tag"`
	MsgId string `json:"msg_id"`
}
