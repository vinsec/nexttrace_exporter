package executor

import (
	"context"
	"fmt"
	"log/slog"
	"os/exec"
	"sync"
	"time"

	"github.com/vinsec/nexttrace_exporter/config"
	"github.com/vinsec/nexttrace_exporter/parser"
)

// ExecutionResult stores the result of a nexttrace execution
type ExecutionResult struct {
	Target    string
	Result    *parser.NextTraceResult
	Duration  time.Duration
	Timestamp time.Time
	Error     error
	Status    string // "success", "error", "timeout"
}

// Executor manages the execution of nexttrace commands for multiple targets
type Executor struct {
	binaryPath      string
	timeout         time.Duration
	results         map[string]*ExecutionResult
	resultsMutex    sync.RWMutex
	targets         []config.Target
	cancelFuncs     map[string]context.CancelFunc
	cancelFuncMutex sync.Mutex
	logger          *slog.Logger
}

// NewExecutor creates a new Executor instance
func NewExecutor(binaryPath string, timeout time.Duration, logger *slog.Logger) *Executor {
	return &Executor{
		binaryPath:  binaryPath,
		timeout:     timeout,
		results:     make(map[string]*ExecutionResult),
		cancelFuncs: make(map[string]context.CancelFunc),
		logger:      logger,
	}
}

// Start begins executing nexttrace for all configured targets
func (e *Executor) Start(ctx context.Context, targets []config.Target) {
	e.targets = targets

	for _, target := range targets {
		go e.runTargetLoop(ctx, target)
	}
}

// Stop gracefully stops all running executions
func (e *Executor) Stop() {
	e.cancelFuncMutex.Lock()
	defer e.cancelFuncMutex.Unlock()

	for _, cancel := range e.cancelFuncs {
		cancel()
	}
	e.cancelFuncs = make(map[string]context.CancelFunc)
}

// runTargetLoop runs nexttrace for a single target in a loop
func (e *Executor) runTargetLoop(ctx context.Context, target config.Target) {
	ticker := time.NewTicker(target.Interval)
	defer ticker.Stop()

	// Execute immediately on start
	e.executeTarget(ctx, target)

	for {
		select {
		case <-ctx.Done():
			e.logger.Info("Stopping execution loop for target",
				"target", target.Name,
				"host", target.Host)
			return
		case <-ticker.C:
			e.executeTarget(ctx, target)
		}
	}
}

// executeTarget executes nexttrace for a single target
func (e *Executor) executeTarget(parentCtx context.Context, target config.Target) {
	startTime := time.Now()
	e.logger.Info("Starting nexttrace execution",
		"target", target.Name,
		"host", target.Host)

	// Create a context with timeout
	ctx, cancel := context.WithTimeout(parentCtx, e.timeout)
	defer cancel()

	// Store cancel function
	e.cancelFuncMutex.Lock()
	e.cancelFuncs[target.Name] = cancel
	e.cancelFuncMutex.Unlock()

	// Build command arguments with -j flag for JSON output, --no-color to disable ANSI color codes, and -M to disable map upload
	args := []string{"-j", "--no-color", "-M", target.Host}
	if target.MaxHops > 0 {
		args = append(args, "--max-hops", fmt.Sprintf("%d", target.MaxHops))
	}

	// Execute command
	cmd := exec.CommandContext(ctx, e.binaryPath, args...)
	output, err := cmd.CombinedOutput()
	duration := time.Since(startTime)

	result := &ExecutionResult{
		Target:    target.Name,
		Duration:  duration,
		Timestamp: time.Now(),
	}

	if ctx.Err() == context.DeadlineExceeded {
		result.Status = "timeout"
		result.Error = fmt.Errorf("execution timeout after %v", e.timeout)
		e.logger.Error("NextTrace execution timeout",
			"target", target.Name,
			"host", target.Host,
			"duration", duration)
	} else if err != nil {
		result.Status = "error"
		result.Error = fmt.Errorf("execution failed: %w", err)
		e.logger.Error("NextTrace execution failed",
			"target", target.Name,
			"host", target.Host,
			"error", err,
			"output", string(output))
	} else {
		// Parse the output
		parsed, parseErr := parser.ParseNextTraceOutput(output)
		if parseErr != nil {
			result.Status = "error"
			result.Error = fmt.Errorf("failed to parse output: %w", parseErr)
			e.logger.Error("Failed to parse nexttrace output",
				"target", target.Name,
				"host", target.Host,
				"error", parseErr,
				"output", string(output))
		} else {
			result.Status = "success"
			result.Result = parsed
			e.logger.Info("NextTrace execution completed successfully",
				"target", target.Name,
				"host", target.Host,
				"duration", duration,
				"hops", len(parsed.Hops))
		}
	}

	// Store the result
	e.resultsMutex.Lock()
	e.results[target.Name] = result
	e.resultsMutex.Unlock()

	// Clean up cancel function
	e.cancelFuncMutex.Lock()
	delete(e.cancelFuncs, target.Name)
	e.cancelFuncMutex.Unlock()
}

// GetResult returns the latest result for a target
func (e *Executor) GetResult(targetName string) (*ExecutionResult, bool) {
	e.resultsMutex.RLock()
	defer e.resultsMutex.RUnlock()

	result, exists := e.results[targetName]
	return result, exists
}

// GetAllResults returns all current results
func (e *Executor) GetAllResults() map[string]*ExecutionResult {
	e.resultsMutex.RLock()
	defer e.resultsMutex.RUnlock()

	// Create a copy to avoid race conditions
	results := make(map[string]*ExecutionResult, len(e.results))
	for k, v := range e.results {
		results[k] = v
	}
	return results
}

// Reload updates the targets and restarts the execution loops
func (e *Executor) Reload(ctx context.Context, targets []config.Target) {
	e.logger.Info("Reloading executor with new targets", "count", len(targets))

	// Stop all current executions
	e.Stop()

	// Clear old results for targets that no longer exist
	newTargetNames := make(map[string]bool)
	for _, target := range targets {
		newTargetNames[target.Name] = true
	}

	e.resultsMutex.Lock()
	for name := range e.results {
		if !newTargetNames[name] {
			delete(e.results, name)
		}
	}
	e.resultsMutex.Unlock()

	// Start new executions
	e.Start(ctx, targets)
}
