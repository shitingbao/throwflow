package tool

import (
	"context"
	"crypto/md5"
	"fmt"
	"github.com/go-kratos/kratos/v2/transport"
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

func GetToken() string {
	timestamp := strconv.Itoa(time.Now().Nanosecond())
	hash := md5.New()
	hash.Write([]byte(timestamp))

	return fmt.Sprintf("%x", hash.Sum(nil))
}

func TimeToString(layout string, time time.Time) string {
	return time.Format(layout)
}

func StringToTime(layout, value string) (time.Time, error) {
	return time.ParseInLocation(layout, value, time.Local)
}

func GetGRPCErrorInfo(err error) string {
	grpcErr, _ := status.FromError(err)

	return grpcErr.Message()
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
