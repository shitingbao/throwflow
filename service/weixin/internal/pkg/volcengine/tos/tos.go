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
