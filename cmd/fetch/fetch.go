package cmd

import (
	"context"
	"fmt"
	"github.com/governet/go-fec-fetcher/pkg/downloader"
	"github.com/spf13/cobra"
	"log"
)

const fecBulkDataUrlRoot = "https://fec.gov/files/bulk-downloads/2008/cm08.zip"

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
	fmt.Println("THIS IS ANNOYING")
	downloader.DownloadFile(fecBulkDataUrlRoot, fd.filename)
	return nil
}
