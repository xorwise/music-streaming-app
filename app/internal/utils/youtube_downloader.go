package utils

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func FindAndSaveTrack(ctx context.Context, trackCh chan error, title string, artist string) (string, error) {
	cmd := exec.Command("yt-dlp", fmt.Sprintf("ytsearch:\"%s - %s\"", artist, title), "--skip-download", "--print", "%(title)s $ %(id)s")

	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()

	if err != nil {
		return "", err
	}

	result := strings.TrimSpace(out.String())
	if result == "" {
		return "", fmt.Errorf("track not found")
	}

	parts := strings.Split(result, "$")
	if len(parts) < 2 {
		return "", fmt.Errorf("track not found")
	}

	if strings.Contains(strings.ToLower(parts[0]), strings.ToLower(title)) {
		path := fmt.Sprintf("media/%s.mp3", strings.TrimSpace(parts[1]))
		go downloadTrack(trackCh, parts[1], path)
		return strings.Replace(path, ".mp3", ".m3u8", 1), nil
	}

	return "", fmt.Errorf("track not found")
}

func downloadTrack(trackCh chan error, trackID string, output string) {
	downloadCmd := exec.Command("yt-dlp", fmt.Sprintf("https://www.youtube.com/watch?v=%s", strings.TrimSpace(trackID)), "--extract-audio", "--audio-format", "mp3", "--audio-quality", "0", "-o", output)
	downloadCmd.Stdout = os.Stdout
	downloadCmd.Stderr = os.Stderr
	err := downloadCmd.Run()
	if err != nil {
		trackCh <- err
	}
	filename := strings.Replace(output, ".mp3", "", 1)
	convertToM3U8Cmd := exec.Command("ffmpeg", "-i", output, "-hls_time", "10", "-hls_playlist_type", "vod", "-hls_segment_filename", filename+"_%03d.ts", strings.Replace(output, ".mp3", ".m3u8", 1))
	convertToM3U8Cmd.Stdout = os.Stdout
	convertToM3U8Cmd.Stderr = os.Stderr
	err = convertToM3U8Cmd.Run()
	if err != nil {
		trackCh <- err
	}
	err = os.Remove(output)
	if err != nil {
		trackCh <- err
	}

	trackCh <- nil
}
