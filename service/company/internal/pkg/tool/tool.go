package tool

import (
	"context"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/shopspring/decimal"
	"golang.org/x/crypto/sha3"
	"google.golang.org/grpc/status"
	"strconv"
	"time"
)

func GetClientIp(ctx context.Context) string {
	remoteIp := "127.0.0.1"

	if tr, err := transport.FromServerContext(ctx); err {
		if ip := tr.RequestHeader().Get("X-Real-IP"); ip != "" {
			remoteIp = ip
		} else if ip = tr.RequestHeader().Get("X-Forwarded-For"); ip != "" {
			remoteIp = ip
		} else if ip = tr.RequestHeader().Get("X-RemoteAddr"); ip != "" {
			remoteIp = ip
		}

		if remoteIp == "::1" {
			remoteIp = "127.0.0.1"
		}
	}

	return remoteIp
}

func TimeToString(layout string, time time.Time) string {
	return time.Format(layout)
}

func StringToTime(layout, value string) (time.Time, error) {
	return time.ParseInLocation(layout, value, time.Local)
}

func GetToken() string {
	timestamp := strconv.Itoa(time.Now().Nanosecond())
	hash := md5.New()
	hash.Write([]byte(timestamp))

	return fmt.Sprintf("%x", hash.Sum(nil))
}

func GetDays(year int, month int) (days int) {
	if month != 2 {
		if month == 4 || month == 6 || month == 9 || month == 11 {
			days = 30

		} else {
			days = 31
		}
	} else {
		if ((year%4) == 0 && (year%100) != 0) || (year%400) == 0 {
			days = 29
		} else {
			days = 28
		}
	}

	return
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

func RemoveDuplicateString(slice []string) []string {
	result := []string{}

	for i := 0; i < len(slice); i++ {
		duplicate := true

		for j := 0; j < len(result); j++ {
			if slice[i] == result[j] {
				duplicate = false

				break
			}
		}

		if duplicate {
			result = append(result, slice[i])
		}
	}

	return result
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

func GetGRPCErrorInfo(err error) string {
	grpcErr, _ := status.FromError(err)

	return grpcErr.Message()
}
