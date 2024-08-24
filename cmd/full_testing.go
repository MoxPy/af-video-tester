package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"sync"
)

var (
	fullTestRtmpUrlFlag  string
	fullTestHlsUrlFlag   string
	fullTestDurationFlag int
	fullTestVlcPathFlag  string
)

var fullTestCmd = &cobra.Command{
	Use:     "full-test",
	Short:   "Test both RTMP and HLS stream simultaneously",
	Long:    "Test both RTMP and HLS stream simultaneously, provide your url, duration for your test and your VLC path.",
	Example: `afvt full-test --rtmpurl rtmp://example.com/stream --hlsurl https://example.com/stream.m3u8 --duration 15 --vlc /usr/bin/vlc`,
	Run: func(cmd *cobra.Command, args []string) {
		if fullTestRtmpUrlFlag == "" || fullTestHlsUrlFlag == "" || fullTestDurationFlag == 0 || fullTestVlcPathFlag == "" {
			fmt.Println("Error: All flags (url, duration, and VLC path) must be provided.")
			err := cmd.Usage()
			if err != nil {
				return
			}
			return
		}
		rtmpServerStatus, rtmpPathStatus, firstHlsStatus, secondHlsStatus := FullTesting(fullTestRtmpUrlFlag, fullTestHlsUrlFlag, fullTestDurationFlag, fullTestVlcPathFlag)
		fmt.Printf("RTMP Server Test. Is it working? %t\n", rtmpServerStatus)
		fmt.Printf("RTMP Stream Path Test. Is it working? %t\n", rtmpPathStatus)
		fmt.Printf("HTTP HLS Test. Is it working? %t\n", firstHlsStatus)
		fmt.Printf("VLC HLS Test. Is it working? %t\n", secondHlsStatus)
		fmt.Println("All tests completed")
	},
}

func init() {
	fullTestCmd.Flags().StringVarP(&fullTestRtmpUrlFlag, "rtmpurl", "r", "", "URL of the RTMP stream to test. For example:  rtmp://example.com:1935/streaming")
	fullTestCmd.Flags().StringVarP(&fullTestHlsUrlFlag, "hlsurl", "u", "", "URL of the HLS stream to test. For example: https://yoursite.com/hls/streaming.m3u8")
	fullTestCmd.Flags().IntVarP(&fullTestDurationFlag, "duration", "d", 10, "Duration of the VLC test in seconds. Set duration to 15 for a quick check or set it to 100 or more for a long check.")
	fullTestCmd.Flags().StringVarP(&fullTestVlcPathFlag, "vlc", "v", "/Applications/VLC.app/Contents/MacOS/VLC", "Path to the VLC executable. For example: /Applications/VLC.app/Contents/MacOS/VLC is the default value, you can omit it if you are on MacOS.")

	hlsCmd.MarkFlagRequired("rtmpurl")
	hlsCmd.MarkFlagRequired("hlsurl")
}

// FullTesting performs two simultaneous tests calling both rtmp and hls checks.
func FullTesting(rtmpUrl string, hlsUrl string, d int, VLCPath string) (bool, bool, bool, bool) {
	var (
		serverStatus    bool
		pathStatus      bool
		firstHlsStatus  bool
		secondHlsStatus bool
	)
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		serverStatus, pathStatus = CheckRTMPStatus(rtmpUrl, VLCPath)
		fmt.Println("RTMP Check COMPLETE")
	}()

	go func() {
		defer wg.Done()
		firstHlsStatus = CheckStatus(hlsUrl)
		if firstHlsStatus {
			secondHlsStatus = CheckWithVLC(hlsUrl, d, VLCPath)
		}
		fmt.Println("HLS Check Completed")
	}()

	wg.Wait()
	return serverStatus, pathStatus, firstHlsStatus, secondHlsStatus
}
