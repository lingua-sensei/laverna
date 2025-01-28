package synthesize

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"runtime"

	"github.com/sourcegraph/conc/pool"
)

// BatchRunner handles concurrent processing of synthesize operations
type BatchRunner struct {
	client     *http.Client
	maxWorkers int
	saveFn     func(string, []byte) error
}

// NewBatchRunner creates a new BatchRunner with the given options
func NewBatchRunner(opts ...BatchRunnerOption) *BatchRunner {
	r := &BatchRunner{
		client:     http.DefaultClient,
		maxWorkers: runtime.GOMAXPROCS(0),
		saveFn: func(text string, audio []byte) error {
			return os.WriteFile(text+".mp3", audio, 0600)
		},
	}

	for _, opt := range opts {
		opt(r)
	}

	return r
}

// BatchRunnerOption configures a BatchRunner
type BatchRunnerOption func(*BatchRunner)

// WithClient sets the HTTP client
func WithClient(c *http.Client) BatchRunnerOption {
	return func(r *BatchRunner) {
		r.client = c
	}
}

// WithMaxWorkers sets the maximum number of concurrent workers
func WithMaxWorkers(n int) BatchRunnerOption {
	return func(r *BatchRunner) {
		r.maxWorkers = n
	}
}

// WithSaveFunc sets custom save function
func WithSaveFunc(fn func(string, []byte) error) BatchRunnerOption {
	return func(r *BatchRunner) {
		r.saveFn = fn
	}
}

// Run runs given opts concurrently and stops if encounters an error
func (r *BatchRunner) Run(ctx context.Context, opts []Opt) error {
	p := pool.New().WithContext(ctx).WithMaxGoroutines(r.maxWorkers)

	for _, opt := range opts {
		p.Go(func(ctx context.Context) error {
			audio, err := Run(ctx, r.client, opt)
			if err != nil {
				return fmt.Errorf("Run(%v): %w", opt, err)
			}

			if err := r.saveFn(opt.Text, audio); err != nil {
				return fmt.Errorf("%T.SaveFunc(%v): %w", p, opt.Text, err)
			}
			return nil
		})
	}

	if err := p.Wait(); err != nil {
		return fmt.Errorf("%T.Wait(): %w", p, err)
	}
	return nil
}
