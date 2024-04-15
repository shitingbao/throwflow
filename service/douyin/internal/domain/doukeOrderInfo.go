package domain

import (
	"context"
	"douyin/internal/pkg/tool"
	"strconv"
	"time"
)

type DoukeOrderInfo struct {
	Id                  uint64
	UserId              uint64
	OrderId             string
	ProductId           string
	ProductName         string
	ProductImg          string
	PaySuccessTime      time.Time
	SettleTime          *time.Time
	RefundTime          *time.Time
	ConfirmTime         *time.Time
	TotalPayAmount      float32
	PayGoodsAmount      float32
	AfterSalesStatus    int64
	FlowPoint           string
	EstimatedCommission float32
	RealCommission      float32
	ItemNum             uint64
	CreateTime          time.Time
	UpdateTime          time.Time
}

func NewDoukeOrderInfo(ctx context.Context, userId, itemNum uint64, afterSalesStatus int64, totalPayAmount, payGoodsAmount, estimatedCommission, realCommission float32, orderId, productId, productName, productImg, flowPoint string, paySuccessTime time.Time) *DoukeOrderInfo {
	return &DoukeOrderInfo{
		UserId:              userId,
		OrderId:             orderId,
		ProductId:           productId,
		ProductName:         productName,
		ProductImg:          productImg,
		PaySuccessTime:      paySuccessTime,
		TotalPayAmount:      totalPayAmount,
		PayGoodsAmount:      payGoodsAmount,
		AfterSalesStatus:    afterSalesStatus,
		FlowPoint:           flowPoint,
		EstimatedCommission: estimatedCommission,
		RealCommission:      realCommission,
		ItemNum:             itemNum,
	}
}

func (doi *DoukeOrderInfo) SetUserId(ctx context.Context, userId uint64) {
	doi.UserId = userId
}

func (doi *DoukeOrderInfo) SetOrderId(ctx context.Context, orderId string) {
	doi.OrderId = orderId
}

func (doi *DoukeOrderInfo) SetProductId(ctx context.Context, productId string) {
	doi.ProductId = productId
}

func (doi *DoukeOrderInfo) SetProductName(ctx context.Context, productName string) {
	doi.ProductName = productName
}

func (doi *DoukeOrderInfo) SetProductImg(ctx context.Context, productImg string) {
	doi.ProductImg = productImg
}

func (doi *DoukeOrderInfo) SetPaySuccessTime(ctx context.Context, paySuccessTime time.Time) {
	doi.PaySuccessTime = paySuccessTime
}

func (doi *DoukeOrderInfo) SetSettleTime(ctx context.Context, settleTime *time.Time) {
	doi.SettleTime = settleTime
}

func (doi *DoukeOrderInfo) SetRefundTime(ctx context.Context, refundTime *time.Time) {
	doi.RefundTime = refundTime
}

func (doi *DoukeOrderInfo) SetConfirmTime(ctx context.Context, confirmTime *time.Time) {
	doi.ConfirmTime = confirmTime
}

func (doi *DoukeOrderInfo) SetTotalPayAmount(ctx context.Context, totalPayAmount float32) {
	doi.TotalPayAmount = totalPayAmount
}

func (doi *DoukeOrderInfo) GetTotalPayAmount(ctx context.Context) string {
	return strconv.FormatFloat(tool.Decimal(float64(doi.TotalPayAmount), 2), 'f', 2, 64)
}

func (doi *DoukeOrderInfo) SetPayGoodsAmount(ctx context.Context, payGoodsAmount float32) {
	doi.PayGoodsAmount = payGoodsAmount
}

func (doi *DoukeOrderInfo) SetAfterSalesStatus(ctx context.Context, afterSalesStatus int64) {
	doi.AfterSalesStatus = afterSalesStatus
}

func (doi *DoukeOrderInfo) SetFlowPoint(ctx context.Context, flowPoint string) {
	doi.FlowPoint = flowPoint
}

func (doi *DoukeOrderInfo) SetEstimatedCommission(ctx context.Context, estimatedCommission float32) {
	doi.EstimatedCommission = estimatedCommission
}

func (doi *DoukeOrderInfo) GetEstimatedCommission(ctx context.Context) string {
	return strconv.FormatFloat(tool.Decimal(float64(doi.EstimatedCommission*0.75), 2), 'f', 2, 64)
}

func (doi *DoukeOrderInfo) SetRealCommission(ctx context.Context, realCommission float32) {
	doi.RealCommission = realCommission
}

func (doi *DoukeOrderInfo) GetRealCommission(ctx context.Context) string {
	return strconv.FormatFloat(tool.Decimal(float64(doi.RealCommission*0.75), 2), 'f', 2, 64)
}

func (doi *DoukeOrderInfo) SetItemNum(ctx context.Context, itemNum uint64) {
	doi.ItemNum = itemNum
}

func (doi *DoukeOrderInfo) SetCreateTime(ctx context.Context) {
	doi.CreateTime = time.Now()
}

func (doi *DoukeOrderInfo) SetUpdateTime(ctx context.Context) {
	doi.UpdateTime = time.Now()
}
