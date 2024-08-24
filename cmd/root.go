package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "afvt",
	Short: "AF Video Tester",
	Long: `A simple command-line tool for testing RTMP and HLS streaming. AFV utilizes VLC and HTTP requests to handle and 
analyze the streaming status.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(hlsCmd)
	rootCmd.AddCommand(rtmpCmd)
	rootCmd.AddCommand(longRtmpCmd)
	rootCmd.AddCommand(fullTestCmd)
	rootCmd.PersistentFlags().StringP("author", "a", "Manuel Lanzani", "Author name for copyright attribution")
}
