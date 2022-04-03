package s3

import (
	"bytes"
	"fmt"
	"io"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type S3 struct {
	client   s3.S3
	uploader *s3manager.Uploader
}

func New(region string, creds *credentials.Credentials) (*S3, error) {
	s3client := S3{}
	err := s3client.setupClient(region, creds)
	if err != nil {
		return nil, fmt.Errorf("creating new S3 client: %v", err)
	}

	s3client.setupUploader()
	if err != nil {
		return nil, fmt.Errorf("initializing uploader: %v", err)
	}
	return &s3client, nil
}

func (s *S3) ListObjects(bucket string, nextToken *string) ([]*s3.Object, error) {
	log.Printf("listing objects in bucket %s", bucket)
	output := []*s3.Object{}
	input := &s3.ListObjectsV2Input{
		Bucket: &bucket,
	}
	if nextToken != nil {
		log.Printf("nextToken: %s", *nextToken)
		input.ContinuationToken = nextToken
	}

	lo, err := s.client.ListObjectsV2(input)
	if err != nil {
		return nil, fmt.Errorf("listing objects in S3 bucket %s: %v", bucket, err)
	}
	log.Printf("items listed: %d", len(lo.Contents))
	output = append(output, lo.Contents...)
	if *lo.IsTruncated {
		o, err := s.ListObjects(bucket, lo.NextContinuationToken)
		if err != nil {
			return nil, err
		}
		output = append(output, o...)
	}
	return output, nil
}

func (s *S3) PutObject(bucket string, key string, object []byte) error {
	log.Printf("putting object in bucket %s at key %s", bucket, key)
	input := &s3.PutObjectInput{
		Bucket: &bucket,
		Key:    &key,
		Body:   bytes.NewReader(object),
	}
	_, err := s.client.PutObject(input)
	if err != nil {
		return fmt.Errorf("putting object at key %s in bucket %s", key, bucket)
	}
	return nil
}

func (s *S3) UploadObject(bucket string, key string, object io.Reader) (*s3manager.UploadOutput, error) {
	result, err := s.uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
		Body:   object,
	})
	if err != nil {
		return nil, fmt.Errorf("uploading object to key %s in bucket %s: %v", key, bucket, err)
	}
	return result, nil
}

func (s *S3) setupUploader() {
	s.uploader = s3manager.NewUploaderWithClient(&s.client, func(u *s3manager.Uploader) {
		u.PartSize = 5 * 1024 * 1024 // The minimum/default allowed part size is 5MB
		u.Concurrency = 2            // default is 5
	})
}

func (s *S3) setupClient(region string, creds *credentials.Credentials) error {
	sess, err := session.NewSession(
		&aws.Config{
			Region:      aws.String(region),
			Credentials: creds,
		},
	)
	if err != nil {
		return fmt.Errorf("setting up S3 client: %v", err)
	}
	svc := s3.New(sess)
	s.client = *svc
	return nil
}
