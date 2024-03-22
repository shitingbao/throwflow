package tool

import (
	"google.golang.org/grpc/status"
)

func GetGRPCErrorInfo(err error) string {
	grpcErr, _ := status.FromError(err)

	return grpcErr.Message()
}
