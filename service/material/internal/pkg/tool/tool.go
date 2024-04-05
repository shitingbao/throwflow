package tool

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"github.com/shopspring/decimal"
	"golang.org/x/crypto/sha3"
	"time"
)

func TimeToString(layout string, time time.Time) string {
	return time.Format(layout)
}

func GetRandCode(data string) string {
	b := make([]byte, 4)
	_, err := rand.Read(b)

	hash := sha3.New224()
	hash.Write([]byte(data))

	bytes := hash.Sum(nil)

	if err != nil {
		return hex.EncodeToString(bytes)
	} else {
		return hex.EncodeToString(bytes) + fmt.Sprintf("%x", b)
	}
}

func Decimal(num float64, places int32) float64 {
	num, _ = decimal.NewFromFloat(num).Round(places).Float64()

	return num
}
