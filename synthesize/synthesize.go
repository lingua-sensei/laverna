package synthesize

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type Speed int

const (
	NormalSpeed = iota
	SlowerSpeed
	SlowestSpeed
)

type Opts struct {
	Speed  Speed
	Voice  Voice
	Text   string
	Client http.Client
}

// Request will look as below, since it is a form, the key is f.req
// and the URL encoded value is going to be
/*
	[
		[
			[
			"jQ1olc",
			"[
				\"สวัสดีชาวโลก วันนี้เราจะมาพูดคุยกันถึงปัญหาของโลก\",  // Text
				\"th\",                                // Voice
				null,
				null,
				[2]									   // Speed
			]",
			null,
			"generic"
			]
		]
	]
*/
func makeFormData(opts Opts) (string, error) {
	genericOpts := []any{opts.Text, opts.Voice, nil, nil, []Speed{opts.Speed}}
	rawOpts, err := json.Marshal(genericOpts)
	if err != nil {
		return "", fmt.Errorf("json.Marshal(%v): %v", genericOpts, err)
	}

	genericData := [][][]any{
		{
			{"jQ1olc", string(rawOpts), nil, "generic"},
		},
	}
	rawData, err := json.Marshal(genericData)
	if err != nil {
		return "", fmt.Errorf("json.Marshal(%v): %v", genericData, err)
	}

	form := make(url.Values)
	form.Set("f.req", string(rawData))
	return form.Encode(), nil
}

// Response will look as below, this function parses base64 data to MP3 format
/*
	)]}'
	[
		["wrb.fr","jQ1olc","[\"<base 64 data>\"]", null, null, null, "generic"],
		["di", 208],
		["af.httprm",208,"6046482986355911791",35]
	]
*/
func parseAudio(raw []byte) ([]byte, error) {
	lines := strings.Split(string(raw), "\n")

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
		return nil, errors.New("no audio line found")
	}

	var audioLineSubParts [][]any
	if err := json.Unmarshal([]byte(audioLine), &audioLineSubParts); err != nil {
		return nil, fmt.Errorf("json.Unmarshal(%v): %v", audioLine, err)
	}

	audioLineSubPart, ok := audioLineSubParts[0][2].(string)
	if !ok {
		return nil, errors.New("no audio line sub part found")
	}

	var base64EncodedAudio []string
	if err := json.Unmarshal([]byte(audioLineSubPart), &base64EncodedAudio); err != nil {
		return nil, fmt.Errorf("json.Unmarshal(%v): %v", audioLineSubPart, err)
	}

	if len(base64EncodedAudio) == 0 {
		return nil, errors.New("no base64 encoded audio found")
	}

	audio, err := base64.StdEncoding.DecodeString(base64EncodedAudio[0])
	if err != nil {
		return nil, fmt.Errorf("%T.DecodeString(%v): %v", base64.StdEncoding, base64EncodedAudio[0], err)
	}
	return audio, nil
}

const hostname = "https://translate.google.com"

var ErrTextTooLong = errors.New("text must be less than 200 chars")

func Run(opts Opts) ([]byte, error) {
	const url = hostname + "/_/TranslateWebserverUi/data/batchexecute"

	if len(opts.Text) > 200 {
		return nil, ErrTextTooLong
	}

	formData, err := makeFormData(opts)
	if err != nil {
		return nil, fmt.Errorf("makeFormData(%v): %v", opts, err)
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBufferString(formData))
	if err != nil {
		return nil, fmt.Errorf("http.NewRequest(): %v", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Origin", hostname)
	req.Header.Set("Referer", hostname)
	resp, err := opts.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("%T.Do(): %v", opts.Client, err)
	}
	defer resp.Body.Close()

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("io.ReadAll(): %v", err)
	}
	return parseAudio(raw)
}
