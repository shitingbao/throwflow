package domain

import (
	"time"
)

type Industry struct {
	Id           uint64
	IndustryName string
	Status       uint8
	CreateTime   time.Time
	UpdateTime   time.Time
}
