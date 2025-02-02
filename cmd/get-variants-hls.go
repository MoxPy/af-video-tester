package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"io"
	"log"
	"net/http"
	"regexp"
	"strings"
)

var (
	extractVariantsHlsUrlFlag string
)

type StreamInfo struct {
	Bandwidth  string
	Resolution string
	FrameRate  string
	Codecs     string
	Playlist   string
}

var getVariantsAbrCmd = &cobra.Command{
	Use:     "get-variants-hls",
	Aliases: []string{"get-vhls"},
	Short:   "Extract ABR variants from HLS master playlist",
	Long:    "Extract ABR variants from HLS master playlist. Analyzes the provided HLS master playlist URL and lists all available quality variants with their respective bandwidth and resolution details. Each variant includes information about resolution, bandwidth, and codec parameters if available.",
	Example: `afvt get-variants-hls --url https://example.com/stream.m3u8`,
	Run: func(cmd *cobra.Command, args []string) {
		if extractVariantsHlsUrlFlag == "" {
			fmt.Println("Error: All flags (url) must be provided.")
			err := cmd.Usage()
			if err != nil {
				return
			}
			return
		}
		fmt.Printf("URL: %s\n", extractVariantsHlsUrlFlag)

		streams := GetVariantsHLSUrl(extractVariantsHlsUrlFlag)

		if len(streams) == 0 {
			fmt.Println("No ABR variants found in the provided HLS playlist.")
			return
		}

		fmt.Println("Available Variants:")
		fmt.Println("----------------------------------------------------------------------------")
		fmt.Println("Bandwidth  | Resolution  | Codec  | Playlist")
		fmt.Println("----------------------------------------------------------------------------")
		for _, s := range streams {
			fmt.Printf("%-10s | %-10s | %-20s | %s\n", s.Bandwidth, s.Resolution, s.Codecs, s.Playlist)
		}
		fmt.Println("----------------------------------------------------------------------------")

	},
}

func init() {
	getVariantsAbrCmd.Flags().StringVarP(&extractVariantsHlsUrlFlag, "url", "u", "", "URL of the HLS stream to test. For example: https://yoursite.com/hls/streaming.m3u8")

	err := hlsCmd.MarkFlagRequired("url")
	if err != nil {
		return
	}
}

func GetVariantsHLSUrl(url string) []StreamInfo {
	var streams []StreamInfo

	re := regexp.MustCompile(`#EXT-X-STREAM-INF:.*?BANDWIDTH=(\d+)(?:.*?RESOLUTION=(\d+x\d+))?(?:.*?CODECS="([^"]+)")?.*?\n(.*\.m3u8)`)

	resp, err := http.Get(url)
	if err != nil {
		log.Printf("[GetVariantsHLSUrl] Error during HTTP request to %s: %v", url, err)
		return nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("[GetVariantsHLSUrl] .m3u8 not found, HTTP status: %d\n", resp.StatusCode)
		return nil
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("[GetVariantsHLSUrl] Error while reading:", err)
		return nil
	}

	matches := re.FindAllStringSubmatch(string(body), -1)
	for _, match := range matches {
		stream := StreamInfo{
			Bandwidth:  match[1],
			Resolution: "Unknown",
			Codecs:     "Unknown",
			Playlist:   strings.TrimSpace(match[4]),
		}
		if match[2] != "" {
			stream.Resolution = match[2]
		}
		if match[3] != "" {
			stream.Codecs = match[3]
		}
		streams = append(streams, stream)
	}

	return streams
}
