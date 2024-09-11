package tools

import (
	"bytes"
	"context"
	"fmt"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type MinioConfig struct {
	Endpoint        string
	AccessKeyID     string
	SecretAccessKey string
	UseSSL          bool
}

type Minio struct {
	Client *minio.Client
	Config *MinioConfig
}

func (m *Minio) NewClient() error {
	client, err := minio.New(m.Config.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(m.Config.AccessKeyID, m.Config.SecretAccessKey, ""),
		Secure: m.Config.UseSSL,
	})
	if err != nil {
		return err
	}

	if client == nil {
		return fmt.Errorf("minio client is nil")
	}

	m.Client = client
	return nil
}

func (m *Minio) MakeBucket(ctx context.Context, bucketName string) error {
	err := m.Client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
	if err != nil {
		exists, err := m.Client.BucketExists(ctx, bucketName)
		if err != nil {
			return err
		}
		if exists {
			return nil
		}
	}
	return nil
}

func (m *Minio) PutObject(ctx context.Context, content []byte, bucketName, objectName string, opts minio.PutObjectOptions) (minio.UploadInfo, error) {
	info, err := m.Client.PutObject(ctx, bucketName, objectName, bytes.NewReader(content), int64(len(content)), opts)
	if err != nil {
		return minio.UploadInfo{}, err
	}
	return info, nil
}

func (m *Minio) SetBucketPolicy(ctx context.Context, bucketName string) error {
	policy := fmt.Sprintf(`{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Principal":{"AWS":"*"},"Action":["s3:GetObject"],"Resource":["arn:aws:s3:::%s/*"]}]}`, bucketName)
	err := m.Client.SetBucketPolicy(ctx, bucketName, policy)
	if err != nil {
		return err
	}
	return nil
}
