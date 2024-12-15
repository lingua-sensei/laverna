package synthesize

import (
	"testing"
)

func TestRun(t *testing.T) {
	opts := Opts{
		Text:  "สวัสดีชาวโลก วันนี้เราจะมาพูดคุยกันถึงปัญหาของโลก",
		Voice: ThaiVoice,
		Speed: SlowestSpeed,
	}

	audio, err := Run(opts)
	if err != nil {
		t.Fatalf("Synthesize(%v): %v", opts, err)
	}
	if len(audio) == 0 {
		t.Error("audio must not be empty")
	}
}
