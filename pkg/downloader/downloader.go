package downloader

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func DownloadFile(sourceUrl string, destinationFile string) error {
	resp, err := http.Get(sourceUrl)
	if err != nil {
		return fmt.Errorf("error when GETing data: %s", err)
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)
	fmt.Println("status", resp.Status)
	if resp.StatusCode != 200 {
		return fmt.Errorf("error when fetching data: http status %s %d", resp.Status, resp.StatusCode)
	}

	out, err := os.Create(destinationFile)
	if err != nil {
		return fmt.Errorf("error when creating output file: %s", err)
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return fmt.Errorf("error when writing body to output file: %s", err)
}
