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

type SynthesizeOptions struct {
	Slow  bool   `json:"slow"`
	Text  string `json:"text"`
	Voice string `json:"voice"`
}

func createRequestBody(opts SynthesizeOptions) string {
	values := []any{opts.Text, opts.Voice, opts.Slow, "null"}
	valuesJSON, _ := json.Marshal(values)

	data := [][][]any{
		{
			{"jQ1olc", string(valuesJSON), nil, "generic"},
		},
	}
	dataJSON, _ := json.Marshal(data)

	params := url.Values{}
	params.Set("f.req", string(dataJSON))
	params.Set("rpcids", "jQ1olc")
	params.Set("source-path", "/")
	params.Set("bl", "boq_translate-webserver_20241113.05_p0")
	params.Set("hl", "en")
	params.Set("soc-app", "1")
	params.Set("soc-platform", "1")
	params.Set("soc-device", "2")
	params.Set("_reqid", "2651428")
	params.Set("rt", "c")

	return params.Encode()
}

func Synthesize(opts SynthesizeOptions) ([]byte, error) {
	client := &http.Client{}

	body := createRequestBody(opts)

	req, err := http.NewRequest("POST",
		"https://translate.google.com/_/TranslateWebserverUi/data/batchexecute",
		bytes.NewBufferString(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	req.Header.Set("Origin", "https://translate.google.com")
	req.Header.Set("Referer", "https://translate.google.com/")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %v", err)
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %v", err)
	}

	log.Printf("[INFO] Response: %v", string(responseBody))

	return parseResponse(responseBody)
}

func parseResponse(responseBody []byte) ([]byte, error) {
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
	opts := SynthesizeOptions{
		Slow:  false,
		Text:  "สวัสดีโลก",
		Voice: "th",
	}

	audioData, err := Synthesize(opts)
	if err != nil {
		log.Printf("[ERR] Synthesize(%v): %v\n", opts, err)
		return
	}

	const filename = "hello_world_th.mp3"
	if err := os.WriteFile(filename, audioData, 0644); err != nil {
		log.Printf("[ERR] os.WriteFile(%v): %v\n", audioData, err)
		return
	}
	log.Printf("Successfully saved Thai audio to %v\n", filename)
}
