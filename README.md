# AFVT - A command-line tool written in Go for testing RTMP and HLS streaming. You can also analyze and download HLS content.

## Features

A powerful command-line toolkit written in Go for streaming professionals. Test, analyze, and download RTMP and HLS streams with ease.

## ✨ Key Features

- 🔍 **Stream Testing**: Comprehensive RTMP and HLS stream validation
- ⚡ **Parallel Processing**: Concurrent stream testing and downloads using Go routines
- 📊 **HLS Analysis**: Extract and analyze ABR variants from master playlists
- 📥 **Batch Downloads**: Download multiple HLS streams simultaneously
- 🎮 **VLC Integration**: Seamless integration with VLC for stream verification
- 🔄 **Format Support**: Full RTMP and HLS protocol support

## Usage

### RTMP Test

This command tests a single RTMP stream.
During the RTMP check, the tool performs the following actions:

    IP Address Check: The RTMP check initially verifies the IP address and port of the RTMP server. For example, if you provide an RTMP URL like rtmp://1.1.1.1:1935/streampath, the tool will only check the IP address and port (1.1.1.1:1935).

    Stream Path Verification: The specific stream path (streampath in this example) will be verified separately using VLC. This ensures that while the server connection is established, the actual streaming content is validated in a subsequent step.

```bash 
afvt rtmp --url rtmp://example.com:1935/streaming
```
    --url: URL of the RTMP stream to test.
    --vlc, -v: Path to the VLC executable. For example: /Applications/VLC.app/Contents/MacOS/VLC is the default value, you can omit it if you are on MacOS.

### HLS Test

This command tests a single HLS stream.
During the HLS test, the tool follows these steps:

    HTTP Request: The tool first uses the CheckStatus() function to send an HTTP request to the provided HLS URL. This function verifies if the HLS URL is accessible and responds correctly.

    Stream Verification: If CheckStatus() returns a positive response, indicating that the URL is valid and accessible, the tool then proceeds to call CheckWithVLC(). This function uses VLC to further verify the streaming content and ensure that the HLS stream is working properly.

```bash 
afvt hls --url https://example.com/stream.m3u8 --duration 20
```
    --url: URL of the HLS stream to test.
    --duration: Duration of the test in seconds. Default is 20 seconds.
    --vlc: Path to the VLC executable. For example: /Applications/VLC.app/Contents/MacOS/VLC is the default value, you can omit it if you are on MacOS.

### Long RTMP Test

This command tests a single RTMP stream for 100s.

```bash 
afvt long-rtmp --url rtmp://example.com:1935/streaming
```
    --url: URL of the RTMP stream to test.
    --vlc: Path to the VLC executable. For example: /Applications/VLC.app/Contents/MacOS/VLC is the default value, you can omit it if you are on MacOS.

### Full Test

This command tests both RTMP and HLS stream. It takes advantage of Go's goroutines to open and test multiple streams concurrently.

```bash 
afvt full-test --rtmpurl rtmp://localhost:1935/streaming --hlsurl http://localhost:8080/hls/streaming.m3u8
```
    --rtmpurl: URL of the RTMP stream to test.
    --hlsurl: URL of the HLS stream to test.
    --duration: Duration of the HLS test in seconds. Default is 20 seconds.
    --vlc: Path to the VLC executable. For example: /Applications/VLC.app/Contents/MacOS/VLC is the default value, you can omit it if you are on MacOS.

### Get Variants

This command extract ABR variants from HLS master playlist. Analyzes the provided HLS master playlist URL and lists all available quality variants with their respective bandwidth and resolution details. Each variant includes information about resolution, bandwidth, and codec parameters if available.
```bash 
afvt get-variants-hls --url https://example.com/stream.m3u8
```

## Download a video from an HLS URL

**Requires ffmpeg**
Download HLS content using ffmpeg. Insert your url and the output file path and choose to open it at the end of the download or not.
```bash 
afvt dw-hls --url http://example.com/playlist.m3u8 --output output.mp4 --open
```
    --url: URL of the file
    --output: Path to save the downloaded file
    --open: If you want to open the video after the download

## Download multiple videos from a HLS URLs

**Requires ffmpeg**
Download multiple HLS content using ffmpeg in parallel. Insert your urls and output file paths.
```bash 
afvt dw-multiple-hls --url "url1,url2,url3" --output "out1.mp4,out2.mp4,out3.mp4" --open
```
    --url: URLs of the file
    --output: Paths to save the downloaded file
    --open: If you want to open the videos after the download

## License

This project is licensed under the Mozilla Public License 2.0. For more details, refer to the LICENSE file in the repository.

## Disclaimer

AFVT is provided as-is, without any warranties or guarantees of any kind, expressed or implied. The use of this application is at your own risk, and the developer disclaims any responsibility for any damages or losses that may arise from its use.

While efforts have been made to ensure the reliability and accuracy of the code, it is essential to review and test thoroughly before deploying in a production environment. The developer is not liable for any consequences, including but not limited to data loss, system failures, or other issues that may occur during the use of AFVT.

Users are encouraged to contribute to the project, report issues, and participate in discussions. However, the developer reserves the right to make changes to the project without prior notice.

For questions, commercial inquiries or additional information, feel free to contact me via [LinkedIn](https://www.linkedin.com/in/manuel-lanzani-59071b251/).
