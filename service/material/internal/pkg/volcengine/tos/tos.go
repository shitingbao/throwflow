package tos

import (
	"context"
	"github.com/volcengine/ve-tos-golang-sdk/v2/tos"
	"io"
)

type Tos struct {
	AccessKey  string
	SecretKey  string
	Endpoint   string
	Region     string
	BucketName string

	client *tos.ClientV2
}

func (t *Tos) NewClient() error {
	client, err := tos.NewClientV2(t.Endpoint, tos.WithRegion(t.Region), tos.WithCredentials(tos.NewStaticCredentials(t.AccessKey, t.SecretKey)))

	if err != nil {
		return err
	}

	t.client = client

	return nil
}

func (t *Tos) GetObject(ctx context.Context, objectKey string) (*tos.GetObjectV2Output, error) {
	output, err := t.client.GetObjectV2(ctx, &tos.GetObjectV2Input{
		Bucket: t.BucketName,
		Key:    objectKey,
	})

	if err != nil {
		return nil, err
	}

	return output, nil
}

func (t *Tos) PutContent(ctx context.Context, fileName string, content io.Reader) (*tos.PutObjectV2Output, error) {
	output, err := t.client.PutObjectV2(ctx, &tos.PutObjectV2Input{
		PutObjectBasicInput: tos.PutObjectBasicInput{
			Bucket: t.BucketName,
			Key:    fileName,
		},
		Content: content,
	})

	if err != nil {
		return nil, err
	}

	return output, nil
}

func (t *Tos) CreateMultipartUpload(ctx context.Context, objectKey string) (*tos.CreateMultipartUploadV2Output, error) {
	createMultipartOutput, err := t.client.CreateMultipartUploadV2(ctx, &tos.CreateMultipartUploadV2Input{
		Bucket: t.BucketName,
		Key:    objectKey,
	})

	if err != nil {
		return nil, err
	}

	return createMultipartOutput, nil
}

func (t *Tos) UploadPart(ctx context.Context, partNumber int, contentLength int64, objectKey, uploadId string, content io.Reader) (*tos.UploadPartV2Output, error) {
	partOutput, err := t.client.UploadPartV2(ctx, &tos.UploadPartV2Input{
		UploadPartBasicInput: tos.UploadPartBasicInput{
			Bucket:     t.BucketName,
			Key:        objectKey,
			UploadID:   uploadId,
			PartNumber: partNumber,
		},
		Content:       content,
		ContentLength: contentLength,
	})

	if err != nil {
		return nil, err
	}

	return partOutput, nil
}

func (t *Tos) ListParts(ctx context.Context, objectKey, uploadId string) (*tos.ListPartsOutput, error) {
	partOutput, err := t.client.ListParts(ctx, &tos.ListPartsInput{
		Bucket:           t.BucketName,
		Key:              objectKey,
		UploadID:         uploadId,
		PartNumberMarker: 0,
	})

	if err != nil {
		return nil, err
	}

	return partOutput, nil
}

func (t *Tos) CompleteMultipartUpload(ctx context.Context, objectKey, uploadId string, parts []tos.UploadedPartV2) (*tos.CompleteMultipartUploadV2Output, error) {
	completeOutput, err := t.client.CompleteMultipartUploadV2(ctx, &tos.CompleteMultipartUploadV2Input{
		Bucket:   t.BucketName,
		Key:      objectKey,
		UploadID: uploadId,
		Parts:    parts,
	})

	if err != nil {
		return nil, err
	}

	return completeOutput, nil
}

func (t *Tos) AbortMultipartUpload(ctx context.Context, objectKey, uploadId string) (*tos.AbortMultipartUploadOutput, error) {
	output, err := t.client.AbortMultipartUpload(ctx, &tos.AbortMultipartUploadInput{
		Bucket:   t.BucketName,
		Key:      objectKey,
		UploadID: uploadId,
	})

	if err != nil {
		return nil, err
	}

	return output, nil
}

func (t *Tos) DeleteObjectV2(ctx context.Context, objectKey string) (*tos.DeleteObjectV2Output, error) {
	output, err := t.client.DeleteObjectV2(ctx, &tos.DeleteObjectV2Input{
		Bucket: t.BucketName,
		Key:    objectKey,
	})

	if err != nil {
		return nil, err
	}

	return output, nil
}
