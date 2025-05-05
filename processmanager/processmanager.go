package processmanager

import (
	"os"
	"os/exec"
	"sync"
)

type ProcessManager struct {
	currentCmd *exec.Cmd
	mu         sync.Mutex
}

func NewProcessManager() *ProcessManager {
	return &ProcessManager{}
}

func (pm *ProcessManager) StartProcess(cmd *exec.Cmd) error {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	if pm.currentCmd != nil {
		if err := pm.currentCmd.Process.Kill(); err != nil {
			return err
		}
	}

	if err := cmd.Start(); err != nil {
		return err
	}

	pm.currentCmd = cmd
	return nil
}

func (pm *ProcessManager) StopCurrent() {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	if pm.currentCmd != nil {
		_ = pm.currentCmd.Process.Kill()
		pm.currentCmd = nil
	}
}
