package data

import (
	"company/internal/biz"
	"context"
	"errors"
	"github.com/go-kratos/kratos/v2/log"
	ctos "github.com/volcengine/ve-tos-golang-sdk/v2/tos"
	"io"
	"time"
)

type companyMaterialRepo struct {
	data *Data
	log  *log.Helper
}

func NewCompanyMaterialRepo(data *Data, logger log.Logger) biz.CompanyMaterialRepo {
	return &companyMaterialRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (cmr *companyMaterialRepo) SaveCacheHash(ctx context.Context, key string, val map[string]string, timeout time.Duration) error {
	_, err := cmr.data.rdb.HMSet(ctx, key, val).Result()

	if err != nil {
		return err
	}

	_, err = cmr.data.rdb.Expire(ctx, key, timeout).Result()

	if err != nil {
		return err
	}

	return nil
}

func (cmr *companyMaterialRepo) UpdateCacheHash(ctx context.Context, key string, val map[string]string) error {
	_, err := cmr.data.rdb.HMSet(ctx, key, val).Result()

	if err != nil {
		return err
	}

	return nil
}

func (cmr *companyMaterialRepo) GetCacheHash(ctx context.Context, key string, field string) (string, error) {
	val, err := cmr.data.rdb.HGet(ctx, key, field).Result()

	if err != nil {
		return "", err
	}

	return val, nil
}

func (cmr *companyMaterialRepo) DeleteCache(ctx context.Context, key string) error {
	if _, err := cmr.data.rdb.Del(ctx, key).Result(); err != nil {
		return err
	}

	return nil
}

func (cmr *companyMaterialRepo) CreateMultipartUpload(ctx context.Context, objectKey string) (*ctos.CreateMultipartUploadV2Output, error) {
	for _, ltos := range cmr.data.toses {
		if ltos.name == "material" {
			createMultipartOutput, err := ltos.tos.CreateMultipartUpload(ctx, objectKey)

			if err != nil {
				return nil, err
			}

			return createMultipartOutput, nil
		}
	}

	return nil, errors.New("tos is not exist")
}

func (cmr *companyMaterialRepo) UploadPart(ctx context.Context, partNumber int, contentLength int64, objectKey, uploadId string, content io.Reader) (*ctos.UploadPartV2Output, error) {
	for _, ltos := range cmr.data.toses {
		if ltos.name == "material" {
			partOutput, err := ltos.tos.UploadPart(ctx, partNumber, contentLength, objectKey, uploadId, content)

			if err != nil {
				return nil, err
			}

			return partOutput, nil
		}
	}

	return nil, errors.New("tos is not exist")
}

func (cmr *companyMaterialRepo) CompleteMultipartUpload(ctx context.Context, fileName, uploadId string, parts []ctos.UploadedPartV2) (*ctos.CompleteMultipartUploadV2Output, error) {
	for _, ltos := range cmr.data.toses {
		if ltos.name == "material" {
			completeOutput, err := ltos.tos.CompleteMultipartUpload(ctx, fileName, uploadId, parts)

			if err != nil {
				return nil, err
			}

			return completeOutput, nil
		}
	}

	return nil, errors.New("tos is not exist")
}

func (cmr *companyMaterialRepo) AbortMultipartUpload(ctx context.Context, objectKey, uploadId string) (*ctos.AbortMultipartUploadOutput, error) {
	for _, ltos := range cmr.data.toses {
		if ltos.name == "material" {
			output, err := ltos.tos.AbortMultipartUpload(ctx, objectKey, uploadId)

			if err != nil {
				return nil, err
			}

			return output, nil
		}
	}

	return nil, errors.New("tos is not exist")
}

func (cmr *companyMaterialRepo) DeleteObjectV2(ctx context.Context, objectKey string) (*ctos.DeleteObjectV2Output, error) {
	for _, ltos := range cmr.data.toses {
		if ltos.name == "material" {
			output, err := ltos.tos.DeleteObjectV2(ctx, objectKey)

			if err != nil {
				return nil, err
			}

			return output, nil
		}
	}

	return nil, errors.New("tos is not exist")
}
