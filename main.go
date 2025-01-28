package main

import (
	"context"
	"flag"
	"log"
	"os"
	"runtime"

	"github.com/lingua-sensei/laverna/synthesize"
)

var (
	filenamePath = flag.String("file", "", "filename path that is used for reading YAML file")
	maxWorkers   = flag.Int("workers", runtime.GOMAXPROCS(0), "maximum number of concurrent downloads")
)

func main() {
	flag.Parse()
	if *filenamePath == "" {
		flag.Usage()
		os.Exit(0)
	}

	raw, err := os.ReadFile(*filenamePath)
	if err != nil {
		log.Fatalf("[ERR] failed to read filename path: %v", err)
	}

	opts, err := synthesize.UnmarshalYAML(raw)
	if err != nil {
		log.Fatalf("[ERR] failed to unmarshal YAML: %v", err)
	}

	runner := synthesize.NewBatchRunner(synthesize.WithMaxWorkers(*maxWorkers))
	if err := runner.Run(context.Background(), opts); err != nil {
		log.Fatalf("[ERR] failed to run batch: %v", err)
	}
}
