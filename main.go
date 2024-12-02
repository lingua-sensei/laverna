package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type Speed int

const (
	Normal Speed = iota
	Slower
	Slowest
)

type SynthesizeOpts struct {
	Speed Speed
	Text  string
	Voice string
}

func makeFormData(opts SynthesizeOpts) string {
	values := []any{opts.Text, opts.Voice, nil, nil, []Speed{opts.Speed}}
	rawValues, _ := json.Marshal(values)

	data := [][][]any{
		{
			{"jQ1olc", string(rawValues), nil, "generic"},
		},
	}
	rawData, _ := json.Marshal(data)

	formValues := url.Values{}
	formValues.Set("f.req", string(rawData))
	return formValues.Encode()
}

const hostname = "https://translate.google.com"

func Synthesize(opts SynthesizeOpts) ([]byte, error) {
	client := &http.Client{}
	formData := makeFormData(opts)
	req, err := http.NewRequest("POST",
		hostname+"/_/TranslateWebserverUi/data/batchexecute",
		bytes.NewBufferString(formData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Origin", hostname)
	req.Header.Set("Referer", hostname)
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %v", err)
	}
	defer resp.Body.Close()

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %v", err)
	}

	log.Printf("[INFO] Response: %v", string(raw))

	return makeAudioBuffer(raw)
}

func makeAudioBuffer(responseBody []byte) ([]byte, error) {
	responseStr := string(responseBody)
	lines := strings.Split(responseStr, "\n")

	// Try to find the line that contains the array with base64 audio data
	var audioLine string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		// Skip empty lines, numbers, or the ")]}'
		if line == "" || line == ")]}''" || (len(line) > 0 && line[0] >= '0' && line[0] <= '9') {
			continue
		}

		if strings.HasPrefix(line, "[") { // Found what looks like a JSON array
			var jsonArray [][]any
			if err := json.Unmarshal([]byte(line), &jsonArray); err != nil {
				continue
			}
			// Check if the array contains the audio data
			if len(jsonArray) > 0 && len(jsonArray[0]) > 2 && jsonArray[0][2] != nil {
				audioLine = line
				break
			}
		}
	}

	if audioLine == "" {
		return nil, errors.New("no audio data found in response")
	}

	var jsonResponse [][]any
	if err := json.Unmarshal([]byte(audioLine), &jsonResponse); err != nil {
		return nil, fmt.Errorf("failed to parse outer array: %v", err)
	}

	audioDataStr, ok := jsonResponse[0][2].(string)
	if !ok {
		return nil, errors.New("invalid audio data format")
	}

	var audioArray []string
	if err := json.Unmarshal([]byte(audioDataStr), &audioArray); err != nil {
		return nil, fmt.Errorf("failed to parse audio array: %v", err)
	}

	if len(audioArray) == 0 {
		return nil, errors.New("no audio data found in array")
	}

	audioData, err := base64.StdEncoding.DecodeString(audioArray[0])
	if err != nil {
		return nil, fmt.Errorf("failed to decode base64 audio data: %v", err)
	}

	return audioData, nil
}

func main() {
	opts := SynthesizeOpts{
		Text:  "สวัสดีชาวโลก วันนี้เราจะมาพูดคุยกันถึงปัญหาของโลก",
		Voice: "th",
		Speed: Slowest,
	}

	audioData, err := Synthesize(opts)
	if err != nil {
		log.Printf("[ERR] Synthesize(%v): %v\n", opts, err)
		return
	}

	const filename = "hello_world_thai_slowest.mp3"
	if err := os.WriteFile(filename, audioData, 0644); err != nil {
		log.Printf("[ERR] os.WriteFile(%v): %v\n", audioData, err)
		return
	}
	log.Printf("Successfully saved audio to %v\n", filename)
}
