package executor

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/cicd/internal/pipeline"
)

// Result is the outcome of a single step.
type StepResult struct {
	Name     string
	Status   string // "ok", "failed", "skipped"
	Output   string
	Duration time.Duration
	Error    error
}

// Result is the outcome of a stage.
type StageResult struct {
	Name    string
	Steps   []StepResult
	Status  string // "ok", "failed"
	Elapsed time.Duration
}

// RunStages executes stages sequentially, steps within a stage in parallel.
func RunStages(ctx context.Context, stages []pipeline.Stage) []StageResult {
	var results []StageResult

	for _, stage := range stages {
		select {
		case <-ctx.Done():
			results = append(results, StageResult{
				Name:  stage.Name,
				Steps: skipRemaining(stage),
				Status: "failed",
			})
			continue
		default:
		}

		sr := runStage(ctx, stage)
		results = append(results, sr)

		if sr.Status == "failed" {
			// skip remaining stages on failure
			for _, s := range stages[len(results):] {
				results = append(results, StageResult{
					Name:  s.Name,
					Steps: skipRemaining(s),
					Status: "skipped",
				})
			}
			break
		}
	}
	return results
}

func runStage(ctx context.Context, stage pipeline.Stage) StageResult {
	stageCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	sr := StageResult{Name: stage.Name, Status: "ok"}
	start := time.Now()
	defer func() { sr.Elapsed = time.Since(start) }()

	var wg sync.WaitGroup
	var mu sync.Mutex
	sr.Steps = make([]StepResult, len(stage.Steps))

	for i, step := range stage.Steps {
		wg.Add(1)
		go func(idx int, s pipeline.Step) {
			defer wg.Done()
			result := runStep(stageCtx, s)
			mu.Lock()
			sr.Steps[idx] = result
			if result.Status == "failed" {
				sr.Status = "failed"
				cancel() // cancel sibling steps
			}
			mu.Unlock()
		}(i, step)
	}

	wg.Wait()
	return sr
}

func runStep(ctx context.Context, step pipeline.Step) StepResult {
	result := StepResult{Name: step.Name, Status: "ok"}
	start := time.Now()
	defer func() { result.Duration = time.Since(start) }()

	timeout := 30 * time.Minute
	if step.Timeout != "" {
		d, err := time.ParseDuration(step.Timeout)
		if err == nil {
			timeout = d
		}
	}

	stepCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	cmd := exec.CommandContext(stepCtx, "sh", "-c", step.Cmd)

	if step.WorkDir != "" {
		cmd.Dir = step.WorkDir
	}

	cmd.Env = append(os.Environ(), step.Env...)

	var stdout, stderr strings.Builder
	cmd.Stdout = io.MultiWriter(os.Stdout, &stdout)
	cmd.Stderr = io.MultiWriter(os.Stderr, &stderr)

	err := cmd.Run()
	result.Output = stdout.String() + stderr.String()

	if err != nil {
		result.Status = "failed"
		if stepCtx.Err() == context.DeadlineExceeded {
			result.Error = fmt.Errorf("timeout after %s", timeout)
		} else if ctx.Err() == context.Canceled {
			result.Error = fmt.Errorf("cancelled because a parallel step failed")
		} else {
			result.Error = err
		}
	}
	return result
}

func skipRemaining(stage pipeline.Stage) []StepResult {
	results := make([]StepResult, len(stage.Steps))
	for i, s := range stage.Steps {
		results[i] = StepResult{Name: s.Name, Status: "skipped"}
	}
	return results
}
