package synthesize

import (
	"errors"
	"github.com/google/go-cmp/cmp"
	"github.com/mrwormhole/errdiff"
	"net/http"
	"os"
	"testing"
)

func TestRun(t *testing.T) {
	c := &http.Client{}
	opts := Opt{
		Text:  "สวัสดีชาวโลก วันนี้เราจะมาพูดคุยกันถึงปัญหาของโลก",
		Voice: ThaiVoice,
		Speed: SlowestSpeed,
	}

	audio, err := Run(c, opts)
	if err != nil {
		t.Fatalf("Synthesize(%v): %v", opts, err)
	}
	if len(audio) == 0 {
		t.Error("audio must not be empty")
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
