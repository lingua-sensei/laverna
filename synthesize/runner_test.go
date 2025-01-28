package synthesize

import (
	"context"
	"errors"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/mrwormhole/errdiff"
)

func TestBatchRunner(t *testing.T) {
	const maxWorkers = 5
	saveErr := errors.New("save error")
	temp := t.TempDir()

	tests := []struct {
		name       string
		opts       []Opt
		saveFn     func(string, []byte) error
		wantErr    error
		wantAudios []string
		ctx        context.Context
	}{
		{
			name: "successful batch run",
			opts: []Opt{
				{Text: "test1", Voice: EnglishVoice},
				{Text: "test2", Voice: EnglishVoice},
			},
			saveFn: func(text string, audio []byte) error {
				return os.WriteFile(filepath.Join(temp, text+".mp3"), audio, 0600)
			},
			wantAudios: []string{
				filepath.Join(temp, "test1.mp3"),
				filepath.Join(temp, "test2.mp3"),
			},
			ctx: context.Background(),
		},
		{
			name: "context cancelled",
			opts: []Opt{
				{Text: "test3", Voice: EnglishVoice},
			},
			saveFn:  func(string, []byte) error { return nil },
			wantErr: context.Canceled,
			ctx: func() context.Context {
				ctx, cancel := context.WithCancel(context.Background())
				cancel() // Cancel immediately
				return ctx
			}(),
		},
		{
			name: "custom save function error",
			opts: []Opt{
				{Text: "test4", Voice: EnglishVoice},
			},
			saveFn: func(string, []byte) error {
				return saveErr
			},
			wantErr: saveErr,
			ctx:     context.Background(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			runner := NewBatchRunner(
				WithClient(&http.Client{}),
				WithMaxWorkers(maxWorkers),
				WithSaveFunc(tt.saveFn),
			)

			err := runner.Run(tt.ctx, tt.opts)
			if diff := errdiff.Check(err, tt.wantErr); diff != "" {
				t.Errorf("%T.Run(): err diff=\n%s", runner, diff)
			}

			for _, wantAudio := range tt.wantAudios {
				info, err := os.Stat(wantAudio)
				if os.IsNotExist(err) {
					t.Errorf("want file(%q) doesn't exist", wantAudio)
					continue
				}
				if info.Size() == 0 {
					t.Errorf("want file(%q) has no bytes", wantAudio)
				}
			}
		})
	}
}
