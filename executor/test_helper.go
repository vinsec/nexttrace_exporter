package executor

import (
	"time"

	"github.com/vinsec/nexttrace_exporter/parser"
)

// SetTestResult is a helper method for testing to inject mock results
// This should only be used in tests
func (e *Executor) SetTestResult(targetName string, result *parser.NextTraceResult, duration time.Duration) {
	e.resultsMutex.Lock()
	defer e.resultsMutex.Unlock()

	e.results[targetName] = &ExecutionResult{
		Target:    targetName,
		Result:    result,
		Duration:  duration,
		Timestamp: time.Now(),
		Status:    "success",
		Error:     nil,
	}
}
