package downloader

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func DownloadFile(sourceUrl string, destinationFile string) {
	resp, err := http.Get(sourceUrl)
	if err != nil {
		fmt.Printf("error when GETing data: %s", err)
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)
	fmt.Println("status", resp.Status)
	if resp.StatusCode != 200 {
		return
	}

	// Create the file
	out, err := os.Create(destinationFile)
	if err != nil {
		fmt.Printf("error when creating output file: %s", err)
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	fmt.Printf("erro when writing body to output file: %s", err)
}