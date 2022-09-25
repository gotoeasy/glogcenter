package backup

import (
	"context"
	"glc/conf"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func UploadMinio(localPathFile string, minioObjectName string) error {
	ctx := context.Background()
	endpoint := conf.GetMinioUrl()
	accessKeyID := conf.GetMinioUser()
	secretAccessKey := conf.GetMinioPassword()
	bucketName := conf.GetMinioBucket()

	// Initialize minio client object.
	minioClient, err := minio.New(endpoint, &minio.Options{Creds: credentials.NewStaticV4(accessKeyID, secretAccessKey, ""), Secure: false})
	if err != nil {
		return err
	}

	// Make a new bucket called mymusic.
	exists, err := minioClient.BucketExists(ctx, bucketName)
	if !(err == nil && exists) {
		err = minioClient.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
		if err != nil {
			return err
		}
	}

	// Upload
	_, err = minioClient.FPutObject(ctx, bucketName, minioObjectName, localPathFile, minio.PutObjectOptions{ContentType: "application/x-tar"})
	if err != nil {
		return err
	}

	return nil
}
