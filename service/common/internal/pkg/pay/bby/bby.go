package bby

import (
	"bytes"
	"common/internal/conf"
	"common/internal/pkg/tool"
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

const (
	payService    = "unified.trade.online"
	divideService = "unified.trade.divide"
	refundService = "unified.trade.refund"
	version       = "2.0"
	charset       = "UTF-8"
	signType      = "MD5"
	channelName   = "yeepay"
	orderType0    = "wechat"
	orderType1    = "alipay"
	orderType2    = "unionpay"
	accessType0   = "unionpayH5"
	accessType1   = "upJs"
	accessType2   = "wechatH5"
	accessType3   = "wechatJs"
	accessType4   = "wechatMiniProgram"
	accessType5   = "wechatApp"
	accessType6   = "alipayH5"
	accessType7   = "alipayJs"
	accessType8   = "alipayMiniProgram"
	accessType9   = "alipayApp"
	isSplitAmount = "1"
)

var orderType = [3]string{orderType0, orderType1, orderType2}
var accessType = [10]string{accessType0, accessType1, accessType2, accessType3, accessType4, accessType5, accessType6, accessType7, accessType8, accessType9}

type DataPackage struct {
	TimeStamp string `json:"timeStamp"`
	Package   string `json:"package"`
	PaySign   string `json:"paySign"`
	AppId     string `json:"appId"`
	SignType  string `json:"signType"`
	NonceStr  string `json:"nonceStr"`
}

type PayReplyData struct {
	Version       string `xml:"version"`
	Charset       string `xml:"charset"`
	SignType      string `xml:"sign_type"`
	Status        string `xml:"status"`
	ResultCode    string `xml:"result_code"`
	MchId         string `xml:"mch_id"`
	NonceStr      string `xml:"nonce_str"`
	ErrCode       string `xml:"err_code"`
	ErrMsg        string `xml:"err_msg"`
	Sign          string `xml:"sign"`
	TransactionId string `xml:"transaction_id"`
	OutTradeNo    string `xml:"out_trade_no"`
	DataPackage   string `xml:"data_package"`
}

type PayRequestData struct {
	Service       string `xml:"service"`
	Version       string `xml:"version"`
	Charset       string `xml:"charset"`
	SignType      string `xml:"sign_type"`
	ChannelName   string `xml:"channel_name"`
	MchId         string `xml:"mch_id"`
	OutTradeNo    string `xml:"out_trade_no"`
	Body          string `xml:"body"`
	TotalFee      uint64 `xml:"total_fee"`
	MchCreateIp   string `xml:"mch_create_ip"`
	OrderType     string `xml:"order_type"`
	NonceStr      string `xml:"nonce_str"`
	AccessType    string `xml:"access_type"`
	BuyerId       string `xml:"buyer_id"`
	NotifyUrl     string `xml:"notify_url"`
	IsSplitAmount string `xml:"is_split_amount"`
	Sign          string `xml:"sign"`
}

type PayAsyncNotificationData struct {
	Version          string `xml:"version"`
	Charset          string `xml:"charset"`
	SignType         string `xml:"sign_type"`
	ResultCode       string `xml:"result_code"`
	Status           string `xml:"status"`
	PayResult        string `xml:"pay_result"`
	MchId            string `xml:"mch_id"`
	NonceStr         string `xml:"nonce_str"`
	Sign             string `xml:"sign"`
	TransactionId    string `xml:"transaction_id"`
	TradeType        string `xml:"trade_type"`
	OutTradeNo       string `xml:"out_trade_no"`
	OutTransactionId string `xml:"out_transaction_id"`
	TotalFee         uint64 `xml:"total_fee"`
	FeeType          string `xml:"fee_type"`
	Body             string `xml:"body"`
	NotifyTime       string `xml:"notify_time"`
	TimeEnd          string `xml:"time_end"`
	ErrCode          string `xml:"err_code"`
	ErrMsg           string `xml:"err_msg"`
}

func Pay(outTradeNo, body, mchCreateIp, nonceStr, openId string, totalFee uint64, conf *conf.Pay_BbyAccount) (*PayReplyData, error) {
	sign := tool.GetMd5("access_type=" + accessType[4] + "&body=" + body + "&buyer_id=" + openId + "&channel_name=" + channelName + "&charset=" + charset + "&is_split_amount=" + isSplitAmount + "&mch_create_ip=" + mchCreateIp + "&mch_id=" + conf.MchId + "&nonce_str=" + nonceStr + "&notify_url=" + conf.CallBackUrl + "&order_type=" + orderType[0] + "&out_trade_no=" + outTradeNo + "&service=" + payService + "&sign_type=" + signType + "&total_fee=" + strconv.FormatUint(totalFee, 10) + "&version=" + version + "&key=" + conf.SecretKey)

	requestData := PayRequestData{
		Service:       payService,
		Version:       version,
		Charset:       charset,
		SignType:      signType,
		ChannelName:   channelName,
		MchId:         conf.MchId,
		OutTradeNo:    outTradeNo,
		Body:          body,
		TotalFee:      totalFee,
		MchCreateIp:   mchCreateIp,
		OrderType:     orderType[0],
		NonceStr:      nonceStr,
		AccessType:    accessType[4],
		BuyerId:       openId,
		NotifyUrl:     conf.CallBackUrl,
		IsSplitAmount: isSplitAmount,
		Sign:          sign,
	}

	bytesData, err := xml.Marshal(requestData)

	if err != nil {
		return nil, errors.New(fmt.Sprintf("failed xml marshal : %v", err.Error()))
	}

	request, err := http.NewRequest("POST", conf.Endpoint, bytes.NewReader(bytesData))

	if err != nil {
		return nil, errors.New(fmt.Sprintf("failed new request : %v", err.Error()))
	}

	request.Header.Set("Content-Type", "application/xml;charset=UTF-8")
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

	var replyData PayReplyData

	if err := xml.Unmarshal(respBytes, &replyData); err != nil {
		return nil, errors.New(fmt.Sprintf("failed xml Unmarshal : %v", err.Error()))
	}

	return &replyData, nil
}

func AsyncNotification(content string, conf *conf.Pay_Bby) (*PayAsyncNotificationData, error) {
	var payAsyncNotificationData PayAsyncNotificationData

	if err := xml.Unmarshal([]byte(content), &payAsyncNotificationData); err != nil {
		return nil, errors.New(fmt.Sprintf("failed xml Unmarshal : %v", err.Error()))
	}

	secretKey := ""

	if conf.DjAccount.MchId == payAsyncNotificationData.MchId {
		secretKey = conf.DjAccount.SecretKey
	} else if conf.DefaultAccount.MchId == payAsyncNotificationData.MchId {
		secretKey = conf.DefaultAccount.SecretKey
	}

	sign := tool.GetMd5("body=" + payAsyncNotificationData.Body + "&charset=" + payAsyncNotificationData.Charset + "&fee_type=" + payAsyncNotificationData.FeeType + "&mch_id=" + payAsyncNotificationData.MchId + "&nonce_str=" + payAsyncNotificationData.NonceStr + "&notify_time=" + payAsyncNotificationData.NotifyTime + "&out_trade_no=" + payAsyncNotificationData.OutTradeNo + "&out_transaction_id=" + payAsyncNotificationData.OutTransactionId + "&pay_result=" + payAsyncNotificationData.PayResult + "&result_code=" + payAsyncNotificationData.ResultCode + "&sign_type=" + payAsyncNotificationData.SignType + "&status=" + payAsyncNotificationData.Status + "&time_end=" + payAsyncNotificationData.TimeEnd + "&total_fee=" + strconv.FormatUint(payAsyncNotificationData.TotalFee, 10) + "&trade_type=" + payAsyncNotificationData.TradeType + "&transaction_id=" + payAsyncNotificationData.TransactionId + "&version=" + payAsyncNotificationData.Version + "&key=" + secretKey)

	if strings.ToUpper(sign) != payAsyncNotificationData.Sign {
		return nil, errors.New("异步通知验签失败")
	}

	return &payAsyncNotificationData, nil
}

type DivideReplyData struct {
	Version       string `xml:"version"`
	Charset       string `xml:"charset"`
	Service       string `xml:"service"`
	SignType      string `xml:"sign_type"`
	Status        string `xml:"status"`
	ResultCode    string `xml:"result_code"`
	DivideState   string `xml:"divide_state"`
	MchId         string `xml:"mch_id"`
	NonceStr      string `xml:"nonce_str"`
	DivideDetail  string `xml:"divide_detail"`
	ErrCode       string `xml:"err_code"`
	ErrMsg        string `xml:"err_msg"`
	Sign          string `xml:"sign"`
	TransactionNo string `xml:"transaction_no"`
	OutTradeNo    string `xml:"out_trade_no"`
	OutDivideNo   string `xml:"out_divide_no"`
	DivideOrderNo string `xml:"divide_order_no"`
}

type DivideRequestData struct {
	Service        string `xml:"service"`
	Version        string `xml:"version"`
	Charset        string `xml:"charset"`
	SignType       string `xml:"sign_type"`
	ChannelName    string `xml:"channel_name"`
	MchId          string `xml:"mch_id"`
	OutTradeNo     string `xml:"out_trade_no"`
	TransactionNo  string `xml:"transaction_no"`
	OutDivideNo    string `xml:"out_divide_no"`
	NonceStr       string `xml:"nonce_str"`
	DivideDetail   string `xml:"divide_detail"`
	IsDivideFinish string `xml:"is_divide_finish"`
	Sign           string `xml:"sign"`
}

func Divide(outTradeNo, transactionNo, outDivideNo, nonceStr, toMchId string, amount uint64, conf *conf.Pay_BbyAccount) (*DivideReplyData, error) {
	sign := tool.GetMd5("channel_name=" + channelName + "&charset=" + charset + "&divide_detail=" + "[{'sub_mch_id':'" + toMchId + "','divide_amount':'" + strconv.FormatUint(amount, 10) + "','sub_divide_no':'" + outDivideNo + "1','body':'结算'}]" + "&is_divide_finish=1&mch_id=" + conf.MchId + "&nonce_str=" + nonceStr + "&out_divide_no=" + outDivideNo + "&out_trade_no=" + outTradeNo + "&service=" + divideService + "&sign_type=" + signType + "&transaction_no=" + transactionNo + "&version=" + version + "&key=" + conf.SecretKey)

	requestData := DivideRequestData{
		Service:        divideService,
		Version:        version,
		Charset:        charset,
		SignType:       signType,
		ChannelName:    channelName,
		MchId:          conf.MchId,
		OutTradeNo:     outTradeNo,
		TransactionNo:  transactionNo,
		OutDivideNo:    outDivideNo,
		NonceStr:       nonceStr,
		DivideDetail:   "[{'sub_mch_id':'" + toMchId + "','divide_amount':'" + strconv.FormatUint(amount, 10) + "','sub_divide_no':'" + outDivideNo + "1','body':'结算'}]",
		IsDivideFinish: "1",
		Sign:           sign,
	}
	fmt.Println(requestData)
	bytesData, err := xml.Marshal(requestData)

	if err != nil {
		return nil, errors.New(fmt.Sprintf("failed xml marshal : %v", err.Error()))
	}

	request, err := http.NewRequest("POST", conf.Endpoint, bytes.NewReader(bytesData))

	if err != nil {
		return nil, errors.New(fmt.Sprintf("failed new request : %v", err.Error()))
	}

	request.Header.Set("Content-Type", "application/xml;charset=UTF-8")
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

	var replyData DivideReplyData

	if err := xml.Unmarshal(respBytes, &replyData); err != nil {
		return nil, errors.New(fmt.Sprintf("failed xml Unmarshal : %v", err.Error()))
	}

	return &replyData, nil
}

type RefundReplyData struct {
	Version       string `xml:"version"`
	Charset       string `xml:"charset"`
	SignType      string `xml:"sign_type"`
	ResultCode    string `xml:"result_code"`
	ErrCode       string `xml:"err_code"`
	ErrMsg        string `xml:"err_msg"`
	MchId         string `xml:"mch_id"`
	NonceStr      string `xml:"nonce_str"`
	Sign          string `xml:"sign"`
	TransactionId string `xml:"transaction_id"`
	OutTradeNo    string `xml:"out_trade_no"`
	OutRefundNo   string `xml:"out_refund_no"`
	RefundId      string `xml:"refund_id"`
	RefundFee     uint64 `xml:"refund_fee"`
}

type RefundRequestData struct {
	Service       string `xml:"service"`
	Version       string `xml:"version"`
	Charset       string `xml:"charset"`
	SignType      string `xml:"sign_type"`
	MchId         string `xml:"mch_id"`
	TransactionId string `xml:"transaction_id"`
	OutRefundNo   string `xml:"out_refund_no"`
	TotalFee      uint64 `xml:"total_fee"`
	RefundFee     uint64 `xml:"refund_fee"`
	NonceStr      string `xml:"nonce_str"`
	Sign          string `xml:"sign"`
}

func Refund(outRefundNo, transactionId, nonceStr string, totalFee, refundFee uint64, conf *conf.Pay_BbyAccount) (*RefundReplyData, error) {
	sign := tool.GetMd5("charset=" + charset + "&mch_id=" + conf.MchId + "&nonce_str=" + nonceStr + "&out_refund_no=" + outRefundNo + "&refund_fee=" + strconv.FormatUint(refundFee, 10) + "&service=" + refundService + "&sign_type=" + signType + "&total_fee=" + strconv.FormatUint(totalFee, 10) + "&transaction_id=" + transactionId + "&version=" + version + "&key=" + conf.SecretKey)

	requestData := RefundRequestData{
		Service:       refundService,
		Version:       version,
		Charset:       charset,
		SignType:      signType,
		MchId:         conf.MchId,
		TransactionId: transactionId,
		OutRefundNo:   outRefundNo,
		TotalFee:      totalFee,
		RefundFee:     refundFee,
		NonceStr:      nonceStr,
		Sign:          sign,
	}
	fmt.Println(requestData)
	bytesData, err := xml.Marshal(requestData)

	if err != nil {
		return nil, errors.New(fmt.Sprintf("failed xml marshal : %v", err.Error()))
	}

	request, err := http.NewRequest("POST", conf.Endpoint, bytes.NewReader(bytesData))

	if err != nil {
		return nil, errors.New(fmt.Sprintf("failed new request : %v", err.Error()))
	}

	request.Header.Set("Content-Type", "application/xml;charset=UTF-8")
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
	fmt.Println(respBytes)
	var replyData RefundReplyData

	if err := xml.Unmarshal(respBytes, &replyData); err != nil {
		return nil, errors.New(fmt.Sprintf("failed xml Unmarshal : %v", err.Error()))
	}

	return &replyData, nil
}
