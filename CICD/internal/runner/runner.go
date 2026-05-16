package runner

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/cicd/internal/executor"
	"github.com/cicd/internal/notifier"
	"github.com/cicd/internal/pipeline"
)

// Run loads a pipeline definition, executes it, and sends notifications.
func Run(path string) error {
	p, err := pipeline.Load(path)
	if err != nil {
		return fmt.Errorf("load pipeline: %w", err)
	}

	var notifiers []notifier.Notifier
	if p.Notify.DingTalkURL != "" {
		notifiers = append(notifiers, notifier.NewDingTalk(p.Notify.DingTalkURL))
	}
	if p.Notify.EmailSMTP != "" && p.Notify.EmailTo != "" {
		notifiers = append(notifiers, notifier.NewEmail(
			p.Notify.EmailSMTP, "587", os.Getenv("CICD_EMAIL_USER"),
			os.Getenv("CICD_EMAIL_PASS"), p.Notify.EmailTo,
		))
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)
	go func() {
		<-sigCh
		fmt.Println("\n[runner] received interrupt, cancelling...")
		cancel()
	}()

	fmt.Printf("🚀 Pipeline: %s\n", p.Name)
	startedAt := time.Now()

	stageResults := executor.RunStages(ctx, p.Stages)
	elapsed := time.Since(startedAt)

	status := "ok"
	for _, sr := range stageResults {
		if sr.Status == "failed" {
			status = "failed"
			break
		}
	}

	report := buildReport(p.Name, status, startedAt, elapsed, stageResults)
	printSummary(report)

	if len(notifiers) > 0 {
		fmt.Println("\n📢 Sending notifications...")
		for _, n := range notifiers {
			if err := n.Send(report); err != nil {
				fmt.Printf("  notification failed: %v\n", err)
			}
		}
	}

	if status != "ok" {
		return fmt.Errorf("pipeline %q failed", p.Name)
	}
	return nil
}

func buildReport(name, status string, startedAt time.Time, elapsed time.Duration, stageResults []executor.StageResult) notifier.Report {
	r := notifier.Report{
		Pipeline:  name,
		Status:    status,
		StartedAt: startedAt,
		Elapsed:   elapsed,
	}
	for _, sr := range stageResults {
		srep := notifier.StageReport{
			Name:    sr.Name,
			Status:  sr.Status,
			Elapsed: sr.Elapsed,
		}
		for _, step := range sr.Steps {
			strep := notifier.StepReport{
				Name:     step.Name,
				Status:   step.Status,
				Duration: step.Duration,
			}
			if step.Error != nil {
				strep.Error = step.Error.Error()
			}
			srep.Steps = append(srep.Steps, strep)
		}
		r.Stages = append(r.Stages, srep)
	}
	return r
}

func printSummary(r notifier.Report) {
	fmt.Println()
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Printf("  Pipeline: %s\n", r.Pipeline)
	fmt.Printf("  Status:   %s\n", r.Status)
	fmt.Printf("  Elapsed:  %s\n", r.Elapsed.Truncate(time.Second))
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
}
