package gongmall

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"strings"
	"weixin/internal/pkg/tool"
)

func Sign(paramPattern, appSecret string) string {
	signPattern := paramPattern + "&appSecret=" + appSecret

	return strings.ToUpper(tool.GetMd5(signPattern))
}

func VerifySign(paramPattern, appSecret, sign string) bool {
	signPattern := paramPattern + "&appSecret=" + appSecret

	return sign == strings.ToUpper(tool.GetMd5(signPattern))
}

type CommonResponse struct {
	Success   bool   `json:"success"`
	ErrorCode string `json:"errorCode"`
	ErrorMsg  string `json:"errorMsg"`
	Sign      string `json:"sign"`
}

type AESCipher struct {
	key []byte
}

func NewAESCipher(key string) *AESCipher {
	bkey, _ := base64.StdEncoding.DecodeString(strings.ToUpper(tool.GetMd5(key)))
	return &AESCipher{key: bkey}
}

func (a *AESCipher) Encrypt(raw string) string {
	block, err := aes.NewCipher(a.key)

	if err != nil {
		panic(err)
	}

	rawBytes := []byte(raw)
	blockSize := aes.BlockSize
	rawBytes = PKCS5Padding(rawBytes, blockSize)

	iv := make([]byte, blockSize)
	cipher.NewCBCEncrypter(block, iv).CryptBlocks(rawBytes, rawBytes)

	return base64.StdEncoding.EncodeToString(rawBytes)
}

func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padText...)
}

func RsaEncrypt(publicKey, data string) (string, error) {
	block, _ := pem.Decode([]byte("-----BEGIN PUBLIC KEY-----\n" + publicKey + "\n-----END PUBLIC KEY-----"))

	if block == nil {
		return "", errors.New("failed to decode PEM block")
	}

	pubKey, err := x509.ParsePKIXPublicKey(block.Bytes)

	if err != nil {
		return "", err
	}

	rsaPubKey, ok := pubKey.(*rsa.PublicKey)

	if !ok {
		return "", errors.New("invalid public key type")
	}

	ciphertext, err := rsa.EncryptPKCS1v15(rand.Reader, rsaPubKey, []byte(data))

	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func RsaDecrypt(privateKey, data string) (string, error) {
	block, _ := pem.Decode([]byte("-----BEGIN RSA PRIVATE KEY-----\n" + privateKey + "\n-----END RSA PRIVATE KEY-----"))

	if block == nil {
		return "", errors.New("failed to decode PEM block")
	}

	rprivateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)

	if err != nil {
		return "", errors.New("failed to parse RSA private key")
	}

	bbdata, err := base64.StdEncoding.DecodeString(data)

	if err != nil {
		return "", errors.New("failed to base64 descode")
	}

	decryptedData, err := rsa.DecryptPKCS1v15(nil, rprivateKey, bbdata)

	if err != nil {
		return "", errors.New("failed to decrypt data")
	}

	return string(decryptedData), nil
}

type ContractAsyncNotificationData struct {
	Mobile          string `json:"mobile"`
	Sign            string `json:"sign"`
	Nonce           string `json:"nonce"`
	Identity        string `json:"identity"`
	ContractId      uint64 `json:"contractId"`
	Name            string `json:"name"`
	AppKey          string `json:"appKey"`
	ServiceId       uint64 `json:"serviceId"`
	ContractFileUrl string `json:"contractFileUrl"`
	Timestamp       uint64 `json:"timestamp"`
}

func ContractAsyncNotification(content string) (*ContractAsyncNotificationData, error) {
	var contractAsyncNotificationData ContractAsyncNotificationData

	if err := json.Unmarshal([]byte(content), &contractAsyncNotificationData); err != nil {
		return nil, errors.New(fmt.Sprintf("failed json Unmarshal : %v", err.Error()))
	}

	return &contractAsyncNotificationData, nil
}

type BalanceAsyncNotificationSuccessData struct {
	BankAccount        string  `json:"bankAccount"`
	DateTime           string  `json:"dateTime"`
	Sign               string  `json:"sign"`
	Nonce              string  `json:"nonce"`
	AppKey             string  `json:"appKey"`
	Timestamp          uint64  `json:"timestamp"`
	RequestId          string  `json:"requestId"`
	InnerTradeNo       string  `json:"innerTradeNo"`
	Status             uint8   `json:"status"`
	Name               string  `json:"name"`
	Mobile             string  `json:"mobile"`
	Amount             float32 `json:"amount"`
	CurrentRealWage    float32 `json:"currentRealWage"`
	CurrentTax         float32 `json:"currentTax"`
	CurrentManageFee   float32 `json:"currentManageFee"`
	CurrentAddTax      float32 `json:"currentAddTax"`
	CurrentAddValueTax float32 `json:"currentAddValueTax"`
	Identity           string  `json:"identity"`
	BankName           string  `json:"bankName"`
	Remark             string  `json:"remark,omitempty"`
	CreateTime         string  `json:"createTime"`
	PayTime            string  `json:"payTime"`
}

func BalanceAsyncNotificationSuccess(content string) (*BalanceAsyncNotificationSuccessData, error) {
	var balanceAsyncNotificationSuccessData BalanceAsyncNotificationSuccessData

	if err := json.Unmarshal([]byte(content), &balanceAsyncNotificationSuccessData); err != nil {
		return nil, errors.New(fmt.Sprintf("failed json Unmarshal : %v", err.Error()))
	}

	return &balanceAsyncNotificationSuccessData, nil
}

type BalanceAsyncNotificationFailData struct {
	BankAccount  string  `json:"bankAccount"`
	DateTime     string  `json:"dateTime"`
	InnerTradeNo string  `json:"innerTradeNo"`
	RequestId    string  `json:"requestId"`
	Status       uint8   `json:"status"`
	FailReason   string  `json:"failReason"`
	Sign         string  `json:"sign"`
	Nonce        string  `json:"nonce"`
	AppKey       string  `json:"appKey"`
	Timestamp    uint64  `json:"timestamp"`
	Name         string  `json:"name"`
	Mobile       string  `json:"mobile"`
	Amount       float32 `json:"amount"`
	Identity     string  `json:"identity"`
	CreateTime   string  `json:"createTime"`
	BankName     string  `json:"bankName,omitempty"`
	Remark       string  `json:"remark,omitempty"`
}

func BalanceAsyncNotificationFail(content string) (*BalanceAsyncNotificationFailData, error) {
	var balanceAsyncNotificationFailData BalanceAsyncNotificationFailData

	if err := json.Unmarshal([]byte(content), &balanceAsyncNotificationFailData); err != nil {
		return nil, errors.New(fmt.Sprintf("failed json Unmarshal : %v", err.Error()))
	}

	return &balanceAsyncNotificationFailData, nil
}
