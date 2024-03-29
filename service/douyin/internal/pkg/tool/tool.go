package tool

import (
	"bytes"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/shopspring/decimal"
	"golang.org/x/crypto/sha3"
	"time"
)

func TimeToString(layout string, time time.Time) string {
	return time.Format(layout)
}

func StringToTime(layout, value string) (time.Time, error) {
	return time.ParseInLocation(layout, value, time.Local)
}

func GetMondayOfWeek(layout string, t time.Time) string {
	if t.Weekday() == time.Monday {
		return TimeToString(layout, t)
	} else {
		offset := int(time.Monday - t.Weekday())

		if offset > 0 {
			offset = -6
		}

		return TimeToString(layout, t.AddDate(0, 0, offset))
	}
}

func Marshal(i interface{}) ([]byte, error) {
	buffer := bytes.NewBuffer([]byte{})
	encoder := json.NewEncoder(buffer)
	encoder.SetEscapeHTML(false)
	err := encoder.Encode(i)
	return buffer.Bytes(), err
}

func Decimal(num float64, places int32) float64 {
	num, _ = decimal.NewFromFloat(num).Round(places).Float64()

	return num
}

func RemoveEmptyString(slice []string) []string {
	if len(slice) == 0 {
		return slice
	}

	for i, v := range slice {
		if v == "" {
			slice = append(slice[:i], slice[i+1:]...)
			return RemoveEmptyString(slice)
			break
		}
	}
	return slice
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

func GetMd5(data string) string {
	hash := md5.Sum([]byte(data))

	return hex.EncodeToString(hash[:])
}
