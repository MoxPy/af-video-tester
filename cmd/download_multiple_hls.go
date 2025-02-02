package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"os"
	"os/exec"
	"sync"
)

type DownloadRequest struct {
	URL            string
	OutputFilePath string
	Complete       bool
	Error          error
}

var (
	urlHlsFlags            []string
	outputFileNameHlsFlags []string
	autoOpenMultipleFlag   bool
)

var downloadMultipleHlsCmd = &cobra.Command{
	Use:     "dw-multiple-hls",
	Aliases: []string{"download-multiple-hls"},
	Short:   "Download multiple HLS content using ffmpeg",
	Long:    "Download multiple HLS content using ffmpeg in parallel. Insert your urls and output file paths.",
	Example: `afvt dw-multiple-hls --url "url1,url2,url3" --output "out1.mp4,out2.mp4,out3.mp4" --open`,
	Run: func(cmd *cobra.Command, args []string) {
		if _, err := exec.LookPath("ffmpeg"); err != nil {
			log.Fatal("ffmpeg is not installed. Please install it first: https://ffmpeg.org/download.html")
		}

		if len(urlHlsFlags) != len(outputFileNameHlsFlags) {
			log.Fatal("Number of URLs must match number of output files")
		}

		if len(urlHlsFlags) < 2 {
			log.Fatal("Please provide more than 2 URLs, you can use dw-hls for 1 URL.")
		}

		downloads := make([]DownloadRequest, len(urlHlsFlags))
		for i := range urlHlsFlags {
			downloads[i] = DownloadRequest{
				URL:            urlHlsFlags[i],
				OutputFilePath: outputFileNameHlsFlags[i],
			}
		}

		for _, d := range downloads {
			if _, err := os.Stat(d.OutputFilePath); err == nil {
				fmt.Printf("Warning: File %s already exists. Do you want to overwrite it? (y/n): ", d.OutputFilePath)
				var response string
				_, err := fmt.Scanln(&response)
				if err != nil || (response != "y" && response != "Y") {
					log.Fatal("Operation cancelled by user")
				}
			}
		}

		// wg for goroutines
		var wg sync.WaitGroup
		results := make(chan DownloadRequest, len(downloads))

		fmt.Println("Starting downloads!")
		for i := range downloads {
			wg.Add(1)
			go func(d DownloadRequest) {
				defer wg.Done()
				err := downloadHlsWithFFmpeg(d.URL, d.OutputFilePath)
				results <- DownloadRequest{
					URL:            d.URL,
					OutputFilePath: d.OutputFilePath,
					Complete:       err == nil,
					Error:          err,
				}
			}(downloads[i])
		}

		go func() {
			wg.Wait()
			close(results)
		}()

		successCount := 0
		for result := range results {
			if result.Error != nil {
				fmt.Printf("Failed to download %s: %v\n", result.URL, result.Error)
			} else {
				fmt.Printf("Successfully downloaded to: %s\n", result.OutputFilePath)
				successCount++
				if autoOpenMultipleFlag {
					if err := openFile(result.OutputFilePath); err != nil {
						log.Printf("Warning: Could not open file %s: %v\n", result.OutputFilePath, err)
					}
				}
			}
		}

		fmt.Printf("\nDownload Summary:\n")
		fmt.Printf("Total: %d\n", len(downloads))
		fmt.Printf("Successful: %d\n", successCount)
		fmt.Printf("Failed: %d\n", len(downloads)-successCount)
	},
}

func init() {
	downloadMultipleHlsCmd.Flags().StringSliceVarP(&urlHlsFlags, "url", "u", []string{}, "Comma-separated URLs of the HLS streams")
	downloadMultipleHlsCmd.Flags().StringSliceVarP(&outputFileNameHlsFlags, "output", "o", []string{}, "Comma-separated output file paths")
	downloadMultipleHlsCmd.Flags().BoolVarP(&autoOpenFlag, "open", "O", false, "Automatically open files after download")
	downloadMultipleHlsCmd.MarkFlagRequired("url")
	downloadMultipleHlsCmd.MarkFlagRequired("output")
}
