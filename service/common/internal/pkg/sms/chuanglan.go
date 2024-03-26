package sms

import (
	"bytes"
	"common/internal/conf"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

type ReplyData struct {
	Code     string
	Time     string
	MsgId    string
	ErrorMsg string
}

func Send(conf *conf.Sms, phone, content string) (*ReplyData, error) {
	params := make(map[string]interface{})
	params["account"] = conf.Chuanglan.Account
	params["password"] = conf.Chuanglan.Password
	params["phone"] = phone
	params["msg"] = content

	bytesData, err := json.Marshal(params)

	if err != nil {
		return nil, errors.New(fmt.Sprintf("failed json marshal : %v", err.Error()))
	}

	request, err := http.NewRequest("POST", conf.Chuanglan.Endpoint, bytes.NewReader(bytesData))

	if err != nil {
		return nil, errors.New(fmt.Sprintf("failed new request : %v", err.Error()))
	}

	request.Header.Set("Content-Type", "application/json;charset=UTF-8")
	client := http.Client{}
	resp, err := client.Do(request)

	if err != nil {
		return nil, errors.New(fmt.Sprintf("failed request : %v", err.Error()))
	}

	defer resp.Body.Close()

	respBytes, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, errors.New(fmt.Sprintf("failed file read : %v", err.Error()))
	}

	var replyData ReplyData

	if err := json.Unmarshal(respBytes, &replyData); err != nil {
		return nil, errors.New(fmt.Sprintf("failed json Unmarshal : %v", err.Error()))
	}

	return &replyData, nil
}
