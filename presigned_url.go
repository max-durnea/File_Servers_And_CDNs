package main

/*import(
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"context"
	"time"
)

func generatePresignedURL(s3Client *s3.Client, bucket, key string, expireTime time.Duration) (string, error){
	presignedClient := s3.NewPresignClient(s3Client)
	req, err := presignedClient.PresignGetObject(
        context.TODO(),
        &s3.GetObjectInput{
            Bucket: &bucket,
            Key:    &key,
        },
        s3.WithPresignExpires(expireTime),
    )
	if err != nil {
		return "", err
	}
	return req.URL, nil
}*/