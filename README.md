# AFVT - A command-line tool written in Go for testing RTMP and HLS streaming

## Features

AFVT is a tool that uses a combination of networking techniques and VLC to check if a streaming source is functioning correctly. Here’s how it works:

    Server Check: It uses HTTP requests and a net.Dialer to verify if the streaming servers (RTMP or HLS) are reachable and operational.

    Stream Check: Once the server is confirmed to be up, it launches VLC to test if the streaming feed is actually working and playable. VLC is used to attempt to play the stream and verify its integrity and availability.

You can check an RTMP or HLS streaming source individually, or take advantage of Go's goroutines to open and test multiple streams concurrently using VLC by invoking the full-test command.

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
afvt afvt long-rtmp --url rtmp://example.com:1935/streaming
```
    --url: URL of the RTMP stream to test.
    --vlc: Path to the VLC executable. For example: /Applications/VLC.app/Contents/MacOS/VLC is the default value, you can omit it if you are on MacOS.

### Full Test

This command tests both RTMP and HLS stream. It takes advantage of Go's goroutines to open and test multiple streams concurrently.

```bash 
afvt afvt long-rtmp --url rtmp://example.com:1935/streaming
```
    --rtmpurl: URL of the RTMP stream to test.
    --hlsurl: URL of the HLS stream to test.
    --duration: Duration of the HLS test in seconds. Default is 20 seconds.
    --vlc: Path to the VLC executable. For example: /Applications/VLC.app/Contents/MacOS/VLC is the default value, you can omit it if you are on MacOS.

## License

This project is licensed under the Mozilla Public License 2.0. For more details, refer to the LICENSE file in the repository.

## Disclaimer

AFVT is provided as-is, without any warranties or guarantees of any kind, expressed or implied. The use of this application is at your own risk, and the developer disclaims any responsibility for any damages or losses that may arise from its use.

While efforts have been made to ensure the reliability and accuracy of the code, it is essential to review and test thoroughly before deploying in a production environment. The developer is not liable for any consequences, including but not limited to data loss, system failures, or other issues that may occur during the use of AFVT.

Users are encouraged to contribute to the project, report issues, and participate in discussions. However, the developer reserves the right to make changes to the project without prior notice.

For questions, commercial inquiries or additional information, feel free to contact me via [LinkedIn](https://www.linkedin.com/in/manuel-lanzani-59071b251/).
