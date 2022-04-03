package bucketUrls

import "fmt"

const urlTemplate = "https://%s.s3-%s.amazonaws.com/%s"

type bucketUrl struct {
	bucket string
	region string
}

func New(bucket string, region string) *bucketUrl {
	return &bucketUrl{
		bucket: bucket,
		region: region,
	}
}

func (u *bucketUrl) Url(path string) (url string) {
	return fmt.Sprintf(urlTemplate, u.bucket, u.region, path)
}
