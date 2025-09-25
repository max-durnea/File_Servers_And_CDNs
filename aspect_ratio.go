package main


import (
	"os/exec"
	"encoding/json"
	"bytes"
)

type AspectRatio struct {
	Streams []struct {
		Width int `json:"width"`
		Height int `json:"height"`
	} `json:"streams"`
}
func getVideoAspectRatio(filePath string) (string, error) {
	cmd := exec.Command("ffprobe", "-v", "error", "-print_format", "json", "-show_streams", filePath)
	var out bytes.Buffer
	cmd.Stdout = &out

	if err := cmd.Run(); err != nil {
		return "", err
	}

	var data AspectRatio
	if err := json.Unmarshal(out.Bytes(), &data); err != nil {
		return "", err
	}

	if len(data.Streams) == 0 {
		return "other", nil
	}

	width := data.Streams[0].Width
	height := data.Streams[0].Height

	ratio := float64(width) / float64(height)
	switch {
	case ratio > 1.77 && ratio < 1.78: // ~16:9
		return "16:9",nil
	case ratio > 0.562 && ratio < 0.563: // ~9:16
		return "9:16",nil
	default:
		return "other",nil
	}

}
