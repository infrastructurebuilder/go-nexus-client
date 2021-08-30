package nexus3

import (
	"context"
	"testing"

	minio "github.com/minio/minio-go/v7"
	credentials "github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/stretchr/testify/assert"
)

func TestBlobstoreS3(t *testing.T) {
	bucketName := "s3test"
	bucketLocation := "us-east-1"
	minioGoEndpoint := "localhost:9000"
	minioNexusEndpoint := "http://minio:9000"
	minioAccessKeyID := "minioadmin"
	minioSecretAccessKey := "minioadmin"
	minioUseSSL := false

	err := ensureMinioBucket(bucketName, bucketLocation, minioGoEndpoint, minioUseSSL, minioAccessKeyID, minioSecretAccessKey)
	assert.Nil(t, err)

	client := getTestClient()

	bsName := "test-blobstore-s3"

	bs := &S3BlobStore{
		Name: bsName,
		BucketConfiguration: BlobStoreS3BucketConfiguration{
			Bucket: BlobStoreS3Bucket{
				Name:   bucketName,
				Region: bucketLocation,
			},
			BucketSecurity: &BlobStoreS3BucketSecurity{
				AccessKeyID:     minioAccessKeyID,
				SecretAccessKey: minioSecretAccessKey,
			},
			AdvancedBucketConnection: &BlobStoreS3AdvancedBucketConnection{
				Endpoint:       minioNexusEndpoint,
				ForcePathStyle: true,
			},
		},
	}

	err = client.BlobStore.S3.Create(bs)
	assert.Nil(t, err)

	s3BS, err := client.BlobStore.S3.Get(bs.Name)
	assert.Nil(t, err)
	assert.NotNil(t, s3BS)
	if s3BS != nil {
		assert.NotNil(t, s3BS.BucketConfiguration)
		assert.NotNil(t, s3BS.BucketConfiguration.Bucket)
		assert.NotNil(t, s3BS.BucketConfiguration.BucketSecurity)

		err = client.BlobStore.S3.Delete(bs.Name)
		assert.Nil(t, err)
	}
}

func ensureMinioBucket(bucketName string, bucketLocation string, endpoint string, useSSL bool, accessKeyID string, secretAccessKey string) error {

	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})

	if err != nil {
		return err
	}
	ctx := context.Background()
	exists, err := minioClient.BucketExists(ctx, bucketName)
	if err != nil {
		return err
	}
	if !exists {
		err = minioClient.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{Region: bucketLocation})
	}
	return err
}
