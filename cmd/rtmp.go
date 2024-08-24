package cmd

import (
	"bufio"
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"net"
	"os/exec"
	"strings"
	"time"
)

var (
	rtmpUrlFlag     string
	rtmpVlcPathFlag string
)

var rtmpCmd = &cobra.Command{
	Use:     "rtmp",
	Short:   "Test RTMP stream",
	Aliases: []string{"rtmp-test"},
	Long:    "Test RTMP stream, provide your url and your VLC path.",
	Example: `afvt rtmp --url rtmp://example.com/stream ---vlc /usr/bin/vlc`,
	Run: func(cmd *cobra.Command, args []string) {
		if rtmpUrlFlag == "" || rtmpVlcPathFlag == "" {
			fmt.Println("Error: All flags (url and VLC path) must be provided.")
			err := cmd.Usage()
			if err != nil {
				return
			}
			return
		}
		fmt.Printf("URL: %s\n", rtmpUrlFlag)
		fmt.Printf("VLC Path: %s\n", rtmpVlcPathFlag)
		serverStatus, pathStatus := CheckRTMPStatus(rtmpUrlFlag, rtmpVlcPathFlag)
		fmt.Printf("RTMP Server Status: %t\n", serverStatus)
		fmt.Printf("RTMP Stream Path Status: %t\n", pathStatus)
	},
}

func init() {
	rtmpCmd.Flags().StringVarP(&rtmpUrlFlag, "url", "u", "", "URL of the RTMP stream to test. For example: rtmp://example.com:1935/stream")
	rtmpCmd.Flags().StringVarP(&rtmpVlcPathFlag, "vlc", "v", "/Applications/VLC.app/Contents/MacOS/VLC", "Path to the VLC executable. For example: /Applications/VLC.app/Contents/MacOS/VLC is the default value, you can omit it if you are on MacOS.")

	rtmpCmd.MarkFlagRequired("url")
}

// CheckRTMPStatus checks the status of the RTMP server and the validity of the stream path.
// It first verifies if the RTMP server is reachable by establishing a TCP connection to the host and port specified in the URL.
// If the connection is successful, it proceeds to verify the stream path using VLC media player.
// The VLC media player is launched with the specified URL and the standard output and error streams are monitored.
// The function considers the path status to be valid if VLC outputs the message indicating that the video stream is being processed correctly.
// The function also checks for stream errors in the output from VLC to determine if there was an issue with receiving the stream.
// The function sets the `ServerStatus` and `PathStatus` fields based on the outcome of these checks.
// The operation has a timeout to ensure that VLC is terminated if it takes too long to respond or if there are issues during execution.
func CheckRTMPStatus(url string, VLCPath string) (bool, bool) {
	var serverStatus bool
	var pathStatus bool
	timeout := 35 * time.Second

	urlWithoutProtocol := strings.TrimPrefix(url, "rtmp://")

	slashIndex := strings.Index(urlWithoutProtocol, "/")

	var hostPort, streamPath string
	if slashIndex != -1 {
		hostPort = urlWithoutProtocol[:slashIndex]
		streamPath = urlWithoutProtocol[slashIndex:]
	} else {
		hostPort = urlWithoutProtocol
	}

	// This checks only the ip, not the stream path: ex rtmp://1.1.1.1:1935/streampath
	// Only rtmp://1.1.1.1:1935 will be checked, streampath will be checked later with VLC
	dialer := net.Dialer{Timeout: timeout}
	conn, err := dialer.Dial("tcp", hostPort)

	if err != nil {
		log.Printf("Error connecting to RTMP: %v\n", err)
		return false, false
	}
	serverStatus = true

	defer conn.Close()

	log.Printf("Connected to server RTMP: %s. Stream path %s will be checked next..\n", hostPort, streamPath)
	log.Println("Launching VLC and playing video from RTMP path..")

	// Launch VLC
	//cmd := exec.Command(VLCPath, r.Url, "--verbose", "2")
	cmd := exec.Command(VLCPath, url)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Printf("Error creating stdout output pipe: %v\n", err)
		return serverStatus, false
	}

	stderr, err := cmd.StderrPipe() // VLC uses stderr
	if err != nil {
		log.Printf("Error creating sterr output pipe: %v\n", err)
		return serverStatus, false
	}

	if err = cmd.Start(); err != nil {
		log.Printf("Error starting VLC: %v\n", err)
		return serverStatus, false
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
			log.Println("VLC Output (stderr):", line)
			if strings.Contains(line, "Raising max DPB to 3") {
				go func() {
					select {
					case <-time.After(15 * time.Second):
						pathStatus = true
						log.Println("Test successful, VLC received the signal.")
						done <- true
					}
				}()
			}
			if strings.Contains(line, "stream error") {
				go func() {
					select {
					case <-time.After(5 * time.Second):
						pathStatus = false
						log.Println("Test failed, VLC can't receive the signal.")
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
			return serverStatus, false
		}
		return serverStatus, pathStatus
	case <-time.After(timeout):
		pathStatus = false
		if err = cmd.Process.Kill(); err != nil {
			log.Printf("Error terminating VLC process: %v\n", err)
			return false, false
		} else {
			log.Println("Timeout reached. VLC process has been terminated.")
			return false, false
		}
	}

}
