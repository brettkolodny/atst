package atst_test

import (
	"strings"
	"testing"
	"time"

	"github.com/brettkolodny/atst/pkg/atst"
)

// Helper function to collect outputs with timeout
func collectOutputs(ch chan atst.Output, timeout time.Duration) []atst.Output {
	var outputs []atst.Output
	timeoutCh := time.After(timeout)

	for {
		select {
		case output, ok := <-ch:
			if !ok {
				return outputs
			}
			outputs = append(outputs, output)
		case <-time.After(100 * time.Millisecond):
			// No more outputs for a short time, assume we're done
			return outputs
		case <-timeoutCh:
			// Overall timeout reached
			return outputs
		}
	}
}

func TestStartWithSimpleProgram(t *testing.T) {
	programs := []atst.Program{
		{
			Exec: "echo",
			Args: []string{"hello world"},
		},
	}

	a := atst.Start(programs)

	// Collect outputs
	outputs := collectOutputs(a.Outputs[0], 2*time.Second)

	// Wait for completion
	a.Wait()

	// Check that we got the expected output
	if len(outputs) == 0 {
		t.Fatal("No outputs received")
	}

	output := outputs[0]

	if output.Command != "echo" {
		t.Errorf("Expected command to be 'echo', got '%s'", output.Command)
	}

	if output.Index != 0 {
		t.Errorf("Expected index to be 0, got %d", output.Index)
	}

	if output.Msg != "hello world" {
		t.Errorf("Expected message to be 'hello world', got '%s'", output.Msg)
	}
}

func TestStartWithMultiplePrograms(t *testing.T) {
	programs := []atst.Program{
		{
			Exec: "echo",
			Args: []string{"hello"},
		},
		{
			Exec: "echo",
			Args: []string{"world"},
		},
	}

	a := atst.Start(programs)

	// Collect outputs
	outputs1 := collectOutputs(a.Outputs[0], 2*time.Second)
	outputs2 := collectOutputs(a.Outputs[1], 2*time.Second)

	// Wait for completion
	a.Wait()

	// Check first program outputs
	if len(outputs1) == 0 {
		t.Fatal("No outputs received from first program")
	}

	if outputs1[0].Msg != "hello" {
		t.Errorf("Expected first message to be 'hello', got '%s'", outputs1[0].Msg)
	}

	// Check second program outputs
	if len(outputs2) == 0 {
		t.Fatal("No outputs received from second program")
	}

	if outputs2[0].Msg != "world" {
		t.Errorf("Expected second message to be 'world', got '%s'", outputs2[0].Msg)
	}
}

func TestCaptureStdoutAndStderr(t *testing.T) {
	// Using bash to output to both stdout and stderr
	programs := []atst.Program{
		{
			Exec: "bash",
			Args: []string{"-c", "echo stdout; echo stderr >&2"},
		},
	}

	a := atst.Start(programs)

	// Collect outputs
	outputs := collectOutputs(a.Outputs[0], 2*time.Second)

	// Wait for completion
	a.Wait()

	// Check that we got both stdout and stderr
	var hasStdout, hasStderr bool

	for _, output := range outputs {
		if output.Msg == "stdout" {
			hasStdout = true
		}
		if output.Msg == "stderr" {
			hasStderr = true
		}
	}

	if !hasStdout {
		t.Error("Missing stdout output")
	}

	if !hasStderr {
		t.Error("Missing stderr output")
	}
}

func TestMultipleOutputLines(t *testing.T) {
	programs := []atst.Program{
		{
			Exec: "bash",
			Args: []string{"-c", "echo line1; echo line2; echo line3"},
		},
	}

	a := atst.Start(programs)

	// Collect outputs
	outputs := collectOutputs(a.Outputs[0], 2*time.Second)

	// Wait for completion
	a.Wait()

	// Check that we got all expected lines
	expectedLines := []string{"line1", "line2", "line3"}
	var foundLines int

	for _, expectedLine := range expectedLines {
		for _, output := range outputs {
			if output.Msg == expectedLine {
				foundLines++
				break
			}
		}
	}

	if foundLines != len(expectedLines) {
		t.Errorf("Expected to find %d lines, found %d", len(expectedLines), foundLines)
	}
}

func TestNonExistentProgram(t *testing.T) {
	programs := []atst.Program{
		{
			Exec: "non_existent_program",
			Args: []string{},
		},
	}

	a := atst.Start(programs)

	// Collect outputs
	outputs := collectOutputs(a.Outputs[0], 2*time.Second)

	// Wait for completion
	a.Wait()

	// Check that we got an error
	if len(outputs) == 0 {
		t.Fatal("No outputs received")
	}

	// The error message should mention the command
	var found bool
	for _, output := range outputs {
		if output.Command == "non_existent_program" && output.Msg != "" {
			found = true
			break
		}
	}

	if !found {
		t.Error("Did not receive expected error for non-existent program")
	}
}

func TestProgramWithErrorStatus(t *testing.T) {
	programs := []atst.Program{
		{
			Exec: "bash",
			Args: []string{"-c", "exit 1"},
		},
	}

	a := atst.Start(programs)

	// Collect outputs
	outputs := collectOutputs(a.Outputs[0], 2*time.Second)

	// Wait for completion
	a.Wait()

	// Check that we got an error about the exit status
	var found bool
	for _, output := range outputs {
		if strings.Contains(output.Msg, "exit status 1") {
			found = true
			break
		}
	}

	if !found {
		t.Error("Did not receive expected error with exit status 1")
	}
}

func TestWaitMethod(t *testing.T) {
	programs := []atst.Program{
		{
			Exec: "sleep",
			Args: []string{"0.5"},
		},
	}

	start := time.Now()

	a := atst.Start(programs)
	a.Wait()

	elapsed := time.Since(start)

	// Check that we waited at least 0.5 seconds
	if elapsed < 500*time.Millisecond {
		t.Errorf("Wait() did not block for expected time, elapsed: %v", elapsed)
	}
}

func TestEmptyProgramList(t *testing.T) {
	programs := []atst.Program{}

	a := atst.Start(programs)

	if len(a.Outputs) != 0 {
		t.Errorf("Expected empty outputs slice, got length %d", len(a.Outputs))
	}

	// This should not block
	a.Wait()
}
