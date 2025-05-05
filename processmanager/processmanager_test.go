package processmanager

import (
	"os/exec"
	"testing"
)

func TestStartProcess_Success(t *testing.T) {
	pm := NewProcessManager()

	// Use a no-op command for testing
	cmd := exec.Command("true")

	err := pm.StartProcess(cmd)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if pm.currentCmd == nil {
		t.Error("Expected currentCmd to be set, got nil")
	}
}

func TestStartProcess_ReplaceExisting(t *testing.T) {
	pm := NewProcessManager()

	// Start first process
	firstCmd := exec.Command("true")
	err := pm.StartProcess(firstCmd)
	if err != nil {
		t.Fatalf("Failed to start first process: %v", err)
	}

	// Start second process, should replace the first
	secondCmd := exec.Command("true")
	err = pm.StartProcess(secondCmd)
	if err != nil {
		t.Fatalf("Failed to start second process: %v", err)
	}

	if pm.currentCmd == nil {
		t.Error("Expected currentCmd to be set, got nil")
	}

	// Ensure the first command was killed
	if firstCmd.Process != nil && firstCmd.Process.Pid != 0 {
		t.Error("Expected first command to be killed, but it's still running")
	}
}

func TestStopCurrent(t *testing.T) {
	pm := NewProcessManager()

	// Start a process
	cmd := exec.Command("true")
	err := pm.StartProcess(cmd)
	if err != nil {
		t.Fatalf("Failed to start process: %v", err)
	}

	// Stop the process
	pm.StopCurrent()

	if pm.currentCmd != nil {
		t.Error("Expected currentCmd to be nil after stopping, got non-nil")
	}
}

func TestStartProcess_Error(t *testing.T) {
	pm := NewProcessManager()

	// Use a non-existent command to trigger an error
	cmd := exec.Command("non-existent-command")

	err := pm.StartProcess(cmd)
	if err == nil {
		t.Error("Expected error for non-existent command, got nil")
	}
}
