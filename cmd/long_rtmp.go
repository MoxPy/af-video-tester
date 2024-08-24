package cmd

import (
	"bufio"
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"os/exec"
	"strings"
	"time"
)

var (
	longRtmpUrlFlag     string
	longRtmpVlcPathFlag string
)

var longRtmpCmd = &cobra.Command{
	Use:     "long-rtmp",
	Short:   "Test RTMP stream for 100s",
	Aliases: []string{"long-rtmp-test"},
	Long:    "Test RTMP stream for 100s, provide your url and your VLC path.",
	Example: `afvt long-rtmp --url rtmp://example.com/stream ---vlc /usr/bin/vlc`,
	Run: func(cmd *cobra.Command, args []string) {
		if longRtmpVlcPathFlag == "" || longRtmpUrlFlag == "" {
			fmt.Println("Error: All flags (url and VLC path) must be provided.")
			err := cmd.Usage()
			if err != nil {
				return
			}
			return
		}
		fmt.Printf("URL: %s\n", longRtmpUrlFlag)
		fmt.Printf("VLC Path: %s\n", longRtmpVlcPathFlag)
		pathStatus := LongCheckWithVLC(longRtmpUrlFlag, longRtmpVlcPathFlag)
		fmt.Printf("RTMP Stream Path Status: %t\n", pathStatus)
	},
}

func init() {
	longRtmpCmd.Flags().StringVarP(&longRtmpUrlFlag, "url", "u", "", "URL of the RTMP stream to test. For example: rtmp://example.com:1935/stream")
	longRtmpCmd.Flags().StringVarP(&longRtmpVlcPathFlag, "vlc", "v", "/Applications/VLC.app/Contents/MacOS/VLC", "Path to the VLC executable. For example: /Applications/VLC.app/Contents/MacOS/VLC is the default value, you can omit it if you are on MacOS.")

	longRtmpCmd.MarkFlagRequired("url")
}

// LongCheckWithVLC is used to launch VLC with the provided stream URL. It should only be used if the URL is known to be working.
// If the URL is working, process runs for 100 seconds before being terminated.
func LongCheckWithVLC(url string, VLCPath string) bool {
	var status bool
	cmd := exec.Command(VLCPath, url)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Printf("Error creating stdout output pipe: %v\n", err)
		return false
	}

	stderr, err := cmd.StderrPipe() // VLC uses stderr
	if err != nil {
		log.Printf("Error creating sterr output pipe: %v\n", err)
		return false
	}

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
		fmt.Println("Aborted VLC Stream")
		done <- true
	}()

	go func() {
		for errScanner.Scan() {
			line := errScanner.Text()
			if strings.Contains(line, "Raising max DPB to 3") {
				go func() {
					select {
					case <-time.After(100 * time.Second):
						log.Println("Test successful, VLC received the signal.")
						status = true
						done <- true
					}
				}()
			}
			if strings.Contains(line, "stream error") {
				go func() {
					select {
					case <-time.After(5 * time.Second):
						log.Println("Test failed, VLC can't receive the signal.")
						status = false
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
			return status
		}
		return status
	}
}
