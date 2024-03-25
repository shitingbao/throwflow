package tool

import (
	"context"
	"github.com/go-kratos/kratos/v2/transport"
	"google.golang.org/grpc/status"
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

func GetGRPCErrorInfo(err error) string {
	grpcErr, _ := status.FromError(err)

	return grpcErr.Message()
}

func GetGRPCErrorCode(err error) string {
	grpcErrCode := status.Code(err)

	return grpcErrCode.String()
}

func FormatPhone(phone string) string {
	if len(phone) < 10 {
		return phone
	}

	return phone[:3] + "****" + phone[7:]
}
