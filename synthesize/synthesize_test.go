package synthesize

import (
	"context"
	"errors"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/mrwormhole/errdiff"
)

func TestRun2(t *testing.T) {
	tests := []struct {
		name      string
		client    *http.Client
		opt       Opt
		wantErr   error
		wantBytes int
	}{
		{
			name:   "normal case",
			client: &http.Client{},
			opt: Opt{
				Text:  "สวัสดีชาวโลก วันนี้เราจะมาพูดคุยกันถึงปัญหาของโลก",
				Voice: ThaiVoice,
				Speed: SlowestSpeed,
			},
			wantErr:   nil,
			wantBytes: 66240,
		},
		{
			name:   "text too long",
			client: nil,
			opt: Opt{
				Text:  strings.Repeat("สวัสดีชาวโลก", 200),
				Voice: ThaiVoice,
				Speed: SlowestSpeed,
			},
			wantErr:   ErrTextTooLong,
			wantBytes: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			audio, err := Run(ctx, tt.client, tt.opt)
			if diff := errdiff.Check(err, tt.wantErr); diff != "" {
				t.Errorf("Run(%v): err diff=\n%s", tt.opt, diff)
			}
			if !cmp.Equal(tt.wantBytes, len(audio)) {
				t.Errorf("Run(%v): want bytes(%d) not equal to got bytes(%d)", tt.opt, tt.wantBytes, len(audio))
			}
		})
	}
}

func TestOptsUnmarshalYAML(t *testing.T) {
	tests := []struct {
		name     string
		wantOpts []Opt
		rawYAML  func() []byte
		wantErr  error
	}{
		{
			name: "example YAML",
			rawYAML: func() []byte {
				const filename = "../testdata/synthesize-example.yaml"
				raw, err := os.ReadFile(filename)
				if err != nil {
					t.Fatalf("os.ReadFile(%s): %v", filename, err)
				}
				return raw
			},
			wantOpts: []Opt{
				{
					Speed: NormalSpeed,
					Voice: ThaiVoice,
					Text:  "สวัสดีครับ",
				},
				{
					Speed: SlowerSpeed,
					Voice: EnglishVoice,
					Text:  "Hello there",
				},
				{
					Speed: SlowestSpeed,
					Voice: JapaneseVoice,
					Text:  "こんにちは~",
				},
			},
		},
		{
			name: "empty YAML",
			rawYAML: func() []byte {
				return nil
			},
			wantErr: errors.New("empty yaml"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := UnmarshalYAML(tt.rawYAML())
			if diff := errdiff.Check(err, tt.wantErr); diff != "" {
				t.Errorf("UnmarshalYAML(): err diff=\n%s", diff)
			}

			if diff := cmp.Diff(tt.wantOpts, got); diff != "" {
				t.Errorf("UnmarshalYAML(): opts diff=\n%s", diff)
			}
		})
	}
}
