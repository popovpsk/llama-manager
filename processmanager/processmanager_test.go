package processmanager

import (
	"os/exec"
	"testing"
)

func TestStartProcess_Success(t *testing.T) {
	pm := NewProcessManager()

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

	firstCmd := exec.Command("true")
	err := pm.StartProcess(firstCmd)
	if err != nil {
		t.Fatalf("Failed to start first process: %v", err)
	}

	secondCmd := exec.Command("true")
	err = pm.StartProcess(secondCmd)
	if err != nil {
		t.Fatalf("Failed to start second process: %v", err)
	}

	if pm.currentCmd == nil {
		t.Error("Expected currentCmd to be set, got nil")
	}

	if pm.currentCmd != secondCmd {
		t.Errorf("Expected currentCmd to be secondCmd, got %v", pm.currentCmd)
	}
}

func TestStopCurrent(t *testing.T) {
	pm := NewProcessManager()

	cmd := exec.Command("true")
	err := pm.StartProcess(cmd)
	if err != nil {
		t.Fatalf("Failed to start process: %v", err)
	}

	pm.StopCurrent()

	if pm.currentCmd != nil {
		t.Error("Expected currentCmd to be nil after stopping, got non-nil")
	}
}

func TestStartProcess_Error(t *testing.T) {
	pm := NewProcessManager()

	cmd := exec.Command("non-existent-command")

	err := pm.StartProcess(cmd)
	if err == nil {
		t.Error("Expected error for non-existent command, got nil")
	}
}
