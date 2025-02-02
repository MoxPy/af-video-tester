package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"os"
	"os/exec"
	"runtime"
)

var (
	downloadHlsUrlFlag string
	outputFileNameFlag string
	autoOpenFlag       bool
)

var downloadHlsCmd = &cobra.Command{
	Use:     "dw-hls",
	Aliases: []string{"download-hls"},
	Short:   "Download HLS content using ffmpeg",
	Long:    "Download HLS content using ffmpeg. Insert your url and the output file path and choose to open it at the end of the download or not.",
	Example: `afvt dw-hls --url http://example.com/playlist.m3u8 --output output.mp4 --open`,
	Run: func(cmd *cobra.Command, args []string) {
		if _, err := exec.LookPath("ffmpeg"); err != nil {
			log.Fatal("ffmpeg is not installed. Please install it first: https://ffmpeg.org/download.html")
		}
		if downloadHlsUrlFlag == "" || outputFileNameFlag == "" {
			log.Println("Error: All flags (url, output file path) must be provided.")
		}

		if _, err := os.Stat(outputFileNameFlag); err == nil {
			fmt.Printf("Warning: File %s already exists. Do you want to overwrite it? (y/n): ", outputFileNameFlag)
			var response string
			_, err := fmt.Scanln(&response)
			if err != nil {
				return
			}

			if response != "y" && response != "Y" {
				log.Fatal("Operation cancelled by user")
			}
		} else if !os.IsNotExist(err) {
			log.Fatalf("Error checking file: %s", err)
		}

		err := downloadHlsWithFFmpeg(downloadHlsUrlFlag, outputFileNameFlag)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("Your file has been downloaded!")
		fmt.Printf("You can find it in: %s\n", outputFileNameFlag)
		if autoOpenFlag {
			if err := openFile(outputFileNameFlag); err != nil {
				log.Printf("Warning: Could not open file: %v\n", err)
			}
		}
	},
}

func init() {
	downloadHlsCmd.Flags().StringVarP(&downloadHlsUrlFlag, "url", "u", "", "URL of the HLS stream to test. For example: https://yoursite.com/hls/streaming.m3u8")
	downloadHlsCmd.Flags().StringVarP(&outputFileNameFlag, "output", "o", "", "Output file path")
	downloadHlsCmd.Flags().BoolVarP(&autoOpenFlag, "open", "O", false, "Automatically open the file after download")
	downloadHlsCmd.MarkFlagRequired("url")
}

func openFile(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return fmt.Errorf("file not found: %s", path)
	}
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("cmd", "/c", "start", path)
	case "darwin": // macOS
		cmd = exec.Command("open", path)
	default: // Linux e altri
		cmd = exec.Command("xdg-open", path)
	}

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to open file: %v", err)
	}

	fmt.Printf("Opening %s with default program...\n", path)
	return nil
}

func downloadHlsWithFFmpeg(url string, outputPath string) error {
	cmd := exec.Command("ffmpeg",
		"-i", url, // Input URL
		"-c", "copy", // Copy streams without re-encoding
		"-loglevel", "info", // Show informative messages
		"-stats", // Show progress stats
		"-y",
		outputPath, // Output file
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	log.Printf("Starting download of %s to %s\n", url, outputPath)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("ffmpeg failed: %v", err)
	}

	return nil
}
