package cmd

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/governet/go-fec-fetcher/pkg/bucketUrls"
	"github.com/governet/go-fec-fetcher/pkg/s3"

	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/spf13/cobra"
)

const FecBulkDataUrlRoot = "https://fec.gov/files/bulk-downloads/2008/cm08.zip"
const FecBulkDataS3Bucket = "cg-519a459a-0ea3-42c2-b7bc-fa1143481f74"
const FecBulkDataS3Region = "us-gov-west-1"
const MirrorBucket = "fec-bulk-mirror"
const MirrorRegion = "us-east-1"

type fetchDataOptions struct {
	filename string
}

var fd = &fetchDataOptions{}

var fetchDataCmd = &cobra.Command{
	Use:          "bulk-data",
	Short:        "fetch FEC bulk data",
	Long:         "This command is used to fetch fec bulk data",
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := fd.fetchData(cmd.Context()); err != nil {
			return fmt.Errorf("failed to fetch FEC bulk data: %v", err)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(fetchDataCmd)
	fetchDataCmd.Flags().StringVarP(&fd.filename, "filename", "f", "", "Filename to save bulk data to")
	err := fetchDataCmd.MarkFlagRequired("filename")
	if err != nil {
		log.Fatalf("Error marking flag as required: %v", err)
	}
}

func (fd *fetchDataOptions) fetchData(ctx context.Context) error {
	fecS3Client, err := s3.New(FecBulkDataS3Region, credentials.AnonymousCredentials)
	if err != nil {
		return fmt.Errorf("setting up FEC bulk data s3 client: %v", err)
	}
	u := bucketUrls.New(FecBulkDataS3Bucket, FecBulkDataS3Region)

	fecObjects, err := fecS3Client.ListObjects(FecBulkDataS3Bucket, nil)
	if err != nil {
		return fmt.Errorf("fetching FEC objects: %v", err)
	}

	mirrorCreds := credentials.NewChainCredentials([]credentials.Provider{&credentials.SharedCredentialsProvider{}})
	mirrorS3Client, err := s3.New(MirrorRegion, mirrorCreds)
	if err != nil {
		return fmt.Errorf("setting up mirror bulk data s3 client: %v", err)
	}

	for _, o := range fecObjects {
		k := u.Url(*o.Key)
		fmt.Printf("%v\n", k)
		response, err := http.Get(k)
		if err != nil {
			return fmt.Errorf("reading from url %s: %v", k, err)
		}
		defer response.Body.Close()

		_, err = mirrorS3Client.UploadObject(MirrorBucket, *o.Key, response.Body)
		if err != nil {
			return fmt.Errorf("uploading to mirror bucket %s at key %s: %v", MirrorBucket, *o.Key, err)
		}
	}

	return nil
}
