package main

import (
	"log"
	"os"

	"github.com/lingua-sensei/laverna/synthesize"
)

// New feature, make a CLI that accepts 1 opt or multiple opts via YAML file
// when 1 opt is used via CLI, text should be the filename
// when YAML file used, go with sequential ID name generation but there should be a default start number settable on YAML

func main() {
	opts := synthesize.Opts{
		Text:  "สวัสดีชาวโลก วันนี้เราจะมาพูดคุยกันถึงปัญหาของโลก",
		Voice: synthesize.ThaiVoice,
		Speed: synthesize.NormalSpeed,
	}

	audio, err := synthesize.Run(opts)
	if err != nil {
		log.Printf("[ERR] Synthesize(%v): %v\n", opts, err)
		return
	}

	const filename = "hello_world_thai_slowest.mp3"
	if err := os.WriteFile(filename, audio, 0644); err != nil {
		log.Printf("[ERR] os.WriteFile(%v): %v\n", audio, err)
		return
	}
	log.Printf("[INFO] Successfully saved audio to %v\n", filename)
}
