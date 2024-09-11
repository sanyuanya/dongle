package tools

import (
	"context"
	"testing"

	"github.com/minio/minio-go/v7"
)

func TestMinio(t *testing.T) {

	m := &Minio{
		Config: &MinioConfig{
			Endpoint:        "218.11.1.36:9000",
			AccessKeyID:     "S2TIBM4R4uITLHKhsoy2",
			SecretAccessKey: "HxYsA0TmdvEoy3gGnWA1kq3n2vUkkb0XCMOrj5cZ",
			UseSSL:          false,
		},
	}

	err := m.NewClient()

	if err != nil {
		t.Errorf("Error creating new Minio client: %v", err)
	}

	err = m.MakeBucket(context.Background(), "test-bucket")

	if err != nil {
		t.Errorf("Error creating new Minio bucket: %v", err)
	}

	info, err := m.PutObject(context.Background(), []byte("Hello, Minio!"), "test-bucket", "test-object.txt", minio.PutObjectOptions{
		ContentType:        "text/plain",
		ContentDisposition: "inline",
	})

	if err != nil {
		t.Errorf("Error putting object to Minio: %v", err)
	}

	t.Logf("Successfully created Minio client, bucket and object: %#+v", info)

	t.Logf("Successfully created Minio client, bucket and object")
}

func TestMinio2(t *testing.T) {

	m := &Minio{
		Config: &MinioConfig{
			Endpoint:        "218.11.1.36:9000",
			AccessKeyID:     "EvQqTpffmcfUD91VhnHZ",
			SecretAccessKey: "qAjVfSTMGWS57MQs5z9m4j0Xyr4y17U8dsXOmrmr",
			UseSSL:          false,
		},
	}

	err := m.NewClient()

	if err != nil {
		t.Errorf("Error creating new Minio client: %v", err)
	}

	err = m.MakeBucket(context.Background(), "test4-bucket")

	if err != nil {
		t.Errorf("Error creating new Minio bucket: %v", err)
	}

	info, err := m.PutObject(context.Background(), []byte("Hello, Minio!"), "test4-bucket", "test-object.txt", minio.PutObjectOptions{
		ContentType:        "text/plain",
		ContentDisposition: "inline",
	})

	if err != nil {
		t.Errorf("Error putting object to Minio: %v", err)
	}

	err = m.SetBucketPolicy(context.Background(), "test4-bucket")

	if err != nil {
		t.Errorf("Error setting bucket policy: %v", err)
	}

	t.Logf("Successfully created Minio client, bucket and object: %#+v", info)

	t.Logf("Successfully created Minio client, bucket and object")
}
