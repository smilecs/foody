package data

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

// Function variable for testability
var UploadFileAndGetUrl = uploadFileAndGetUrl

// Real implementation
func uploadFileAndGetUrl(sess *session.Session, bucket, key string, file multipart.File, size int64, contentType string) (string, error) {
	defer file.Close()

	buffer := make([]byte, size)
	_, err := file.Read(buffer)
	if err != nil && err != io.EOF {
		return "", fmt.Errorf("failed to read file: %v", err)
	}

	uploader := s3.New(sess)

	_, err = uploader.PutObject(&s3.PutObjectInput{
		Bucket:        aws.String(bucket),
		Key:           aws.String(key),
		Body:          bytes.NewReader(buffer),
		ContentLength: aws.Int64(size),
		ContentType:   aws.String(contentType),
		ACL:           aws.String("public-read"),
	})
	if err != nil {
		return "", err
	}

	url := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s",
		bucket,
		*sess.Config.Region,
		key,
	)
	return url, err
}
