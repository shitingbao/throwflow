package domain

import (
	"time"
)

type MaterialProduct struct {
	Id           uint64
	ProductId    uint64
	ProductName  string
	VideoLike    uint64
	IndustryId   uint64
	IndustryName string
	CategoryId   uint64
	CategoryName string
	IsHot        uint8
	Videos       uint64
	Awemes       uint64
	Platform     string
	UpdateDay    time.Time
	CreateTime   time.Time
	UpdateTime   time.Time
}
