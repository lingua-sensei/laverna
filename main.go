package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"
	"sync"

	"github.com/lingua-sensei/laverna/synthesize"
)

var (
	filenamePath     = flag.String("file", "", "filename path that is used for reading YAML file")
	maxWorkers       = flag.Int("workers", runtime.GOMAXPROCS(0), "maximum number of concurrent downloads")
	generationNumber = flag.Int("n", 1, "generation number that is used in output filenames")
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

	c := &http.Client{}
	for err := range batchSave(c, *maxWorkers, *generationNumber, opts) {
		log.Printf("[WARN] failed to batch save: %v", err)
	}
}

func batchSave(client *http.Client, workerCount, generationNumber int, opts []synthesize.Opt) <-chan error {
	errChan := make(chan error, len(opts))
	throttle := make(chan struct{}, workerCount)
	var wg sync.WaitGroup
	ctx := context.Background()

	for i := range opts {
		wg.Add(1)
		go func(generationNumber int) {
			defer wg.Done()
			throttle <- struct{}{}
			defer func() {
				<-throttle
			}()

			audio, err := synthesize.Run(ctx, client, opts[i])
			if err != nil {
				errChan <- fmt.Errorf("failed to run opt(%s): %w", opts[i].Text, err)
				return
			}

			filename := fmt.Sprintf("audio_%d.mp3", generationNumber)
			if err := os.WriteFile(filename, audio, 0600); err != nil {
				errChan <- fmt.Errorf("failed to write file(%s): %w", filename, err)
			}
		}(generationNumber + i)
	}

	go func() {
		wg.Wait()
		close(errChan)
	}()
	return errChan
}
