package cmd

import (
	"bufio"
	"fmt"
	"github.com/spf13/cobra"
	"io"
	"log"
	"net/http"
	"os/exec"
	"strings"
	"time"
)

var (
	hlsUrlFlag   string
	durationFlag int
	vlcPathFlag  string
)

var hlsCmd = &cobra.Command{
	Use:     "hls",
	Aliases: []string{"hls-test"},
	Short:   "Test HLS stream",
	Long:    "Test HLS stream, provide your url, duration for your test and your VLC path.",
	Example: `afvt hls --url https://example.com/stream.m3u8 --duration 15 --vlc /usr/bin/vlc`,
	Run: func(cmd *cobra.Command, args []string) {
		if hlsUrlFlag == "" || durationFlag == 0 || vlcPathFlag == "" {
			fmt.Println("Error: All flags (url, duration, and VLC path) must be provided.")
			err := cmd.Usage()
			if err != nil {
				return
			}
			return
		}
		fmt.Printf("URL: %s\n", hlsUrlFlag)
		fmt.Printf("Duration: %s seconds\n", durationFlag)
		fmt.Printf("VLC Path: %s\n", vlcPathFlag)
		firstTestStatus := CheckStatus(hlsUrlFlag)
		fmt.Printf("HTTP HLS Test. Is it working? %t\n", CheckStatus(hlsUrlFlag))
		if firstTestStatus {
			fmt.Printf("VLC HLS Test. Is it working? %t\n", CheckWithVLC(hlsUrlFlag, durationFlag, vlcPathFlag))
		}
	},
}

func init() {
	hlsCmd.Flags().StringVarP(&hlsUrlFlag, "url", "u", "", "URL of the HLS stream to test. For example: https://yoursite.com/hls/streaming.m3u8")
	hlsCmd.Flags().IntVarP(&durationFlag, "duration", "d", 10, "Duration of the VLC test in seconds. Set duration to 15 for a quick check or set it to 100 or more for a long check.")
	hlsCmd.Flags().StringVarP(&vlcPathFlag, "vlc", "v", "/Applications/VLC.app/Contents/MacOS/VLC", "Path to the VLC executable. For example: /Applications/VLC.app/Contents/MacOS/VLC is the default value, you can omit it if you are on MacOS.")

	hlsCmd.MarkFlagRequired("url")
}

// CheckStatus verifies the status of an HLS stream by making an HTTP request to the provided HLS URL.
// It sends a GET request to the URL and checks the HTTP response to determine if the HLS stream is available and active.
// The function first checks for any errors during the HTTP request. If an error occurs, it logs the error and sets the HLS status to `false`.
// If the request is successful, it checks the HTTP status code to ensure the server returned an "OK" status (200).
// The function then reads the response body, which should contain the content of the .m3u8 playlist file.
// It searches the playlist content for the presence of ".ts" (Transport Stream) file references.
// If it finds ".ts" files in the playlist, it indicates that the HLS stream is active and the status is set to `true`.
// If no ".ts" files are found, it indicates that the HLS stream is down and the status is set to `false`.
func CheckStatus(url string) bool {
	resp, err := http.Get(url)

	if err != nil {
		log.Printf("Error during HTTP request: %v\n", err)
		return false
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Printf("Error while closing the response body: %v", err)
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		log.Printf(".m3u8 file not found, HTTP status: %d\n", resp.StatusCode)
		return false
	}

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		log.Printf("Error reading the response: %v\n", err)
		return false
	}

	playlistContent := string(body)

	if strings.Contains(playlistContent, ".ts") {
		log.Printf("HLS Streaming Status: up")
		return true
	}
	log.Printf("HLS Streaming Status: down")
	return false
}

// CheckWithVLC verifies the status of an HLS stream by launching VLC to attempt playback.
// It executes VLC using the provided VLCPath and checks the output to determine if the stream is successfully received.
// Set duration to 15 for a quick check or set it to 100 or more for a long check.
func CheckWithVLC(url string, d int, VLCPath string) bool {
	duration := time.Duration(d) * time.Second
	timeout := duration + 10*time.Second

	cmd := exec.Command(VLCPath, url)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Printf("Error creating stdout output pipe: %v\n", err)
		return false
	}
	defer stdout.Close()
	stderr, err := cmd.StderrPipe() // VLC uses stderr
	if err != nil {
		log.Printf("Error creating stderr pipe: %v\n", err)
		return false
	}
	defer stderr.Close()

	if err = cmd.Start(); err != nil {
		log.Printf("Error starting VLC: %v\n", err)
		return false
	}

	outScanner := bufio.NewScanner(stdout)
	errScanner := bufio.NewScanner(stderr)

	done := make(chan bool, 1)

	go func() {
		for outScanner.Scan() {
			line := outScanner.Text()
			log.Println("VLC Output (stdout):", line)
		}
		done <- true
	}()
	go func() {
		for errScanner.Scan() {
			line := errScanner.Text()
			log.Println("VLC Output (stderr): ", line)
			if strings.Contains(line, "Changing stream format Unknown -> TS") {
				go func() {
					select {
					case <-time.After(duration + 5*time.Second):
						log.Println("Test successful, VLC received the signal.")
						done <- true
					}
				}()
			}
		}
	}()

	select {
	case <-done:
		if err = cmd.Process.Kill(); err != nil {
			log.Printf("Error terminating VLC process: %v\n", err)
			return false
		}
		return true
	case <-time.After(timeout):
		if err = cmd.Process.Kill(); err != nil {
			log.Printf("Error terminating VLC process: %v\n", err)
			return false
		} else {
			log.Println("Timeout reached. VLC process has been terminated.")
			return false
		}
	}
}
