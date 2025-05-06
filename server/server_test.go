package server

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/popovpsk/llama-manager/config"
	"github.com/popovpsk/llama-manager/processmanager"
)

// mockExecCommand is a helper to mock exec.Command
func mockExecCommand(command string, args ...string) *exec.Cmd {
	// cs will be the arguments passed to the test binary.
	// The first argument tells the test binary to run the TestHelperProcess func.
	// The "--" separates the test flags from the command and args we are "mocking".
	cs := []string{"-test.run=TestHelperProcess", "--", command}
	cs = append(cs, args...)
	// os.Args[0] is the path to the current test binary.
	cmd := exec.Command(os.Args[0], cs...)
	// Set this environment variable to signal TestHelperProcess to run.
	cmd.Env = append(os.Environ(), "GO_WANT_HELPER_PROCESS=1")
	return cmd
}

// TestHelperProcess isn't a real test. It's used as a helper process
// for TestHandleShutdown_Success and TestHandleShutdown_CommandError.
// It's executed by the command from mockExecCommand.
func TestHelperProcess(t *testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}
	// Simulate command execution.
	// For Start() error simulation, we could check args or env vars here if needed.
	// For now, successful start is implied by reaching here without os.Exit(1)
	fmt.Fprintf(os.Stdout, "") // Mock successful start
	os.Exit(0)
}

var (
	// To simulate errors from cmd.Start()
	mockStartError      bool
	originalExecCommand = execCommand
)

func TestHandleShutdown(t *testing.T) {
	// Setup: Create a minimal server instance
	// cfg and pm are not directly used by handleShutdown, so nil or empty structs are fine.
	s := NewServer(&config.Config{}, processmanager.NewProcessManager(), "dummy-config.yaml")

	t.Run("Success", func(t *testing.T) {
		// Arrange
		mockStartError = false
		execCommand = func(command string, args ...string) *exec.Cmd {
			cmd := mockExecCommand(command, args...) // Uses TestHelperProcess
			return cmd                               // Return the real *exec.Cmd which will call TestHelperProcess
		}
		defer func() { execCommand = originalExecCommand }() // Restore

		req, err := http.NewRequest("GET", "/shutdown", nil)
		if err != nil {
			t.Fatal(err)
		}
		rr := httptest.NewRecorder()

		// Act
		s.handleShutdown(rr, req)

		// Assert
		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusOK)
		}

		expectedBody := "Shutdown command issued to PC. The system should shut down shortly."
		if !strings.Contains(rr.Body.String(), expectedBody) {
			t.Errorf("handler returned unexpected body: got %v want to contain %v",
				rr.Body.String(), expectedBody)
		}
	})

	t.Run("CommandError", func(t *testing.T) {
		// Arrange
		mockStartError = true // Signal that execCommand should return a failing command
		// Save the current execCommand (which should be originalExecCommand at this point if tests are serial)
		// and restore it afterwards.
		savedExecCommand := execCommand
		execCommand = func(name string, arg ...string) *exec.Cmd {
			if mockStartError {
				// Return a command that is known to fail on Start().
				// Using "/dev/null" or any non-executable path.
				// The actual error message from Start() will be OS-dependent.
				return exec.Command("/dev/null")
			}
			// Fallback to the standard mock if mockStartError is false (though not expected in this test path)
			return mockExecCommand(name, arg...)
		}
		defer func() {
			execCommand = savedExecCommand // Restore to what it was before this subtest
			mockStartError = false
		}()

		req, err := http.NewRequest("GET", "/shutdown", nil)
		if err != nil {
			t.Fatal(err)
		}
		rr := httptest.NewRecorder()

		// Act
		s.handleShutdown(rr, req)

		// Assert
		if status := rr.Code; status != http.StatusInternalServerError {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusInternalServerError)
		}

		// The server formats the error as "Error sending shutdown command: %v".
		// The specific error from exec.Command("/dev/null").Start() is OS-dependent.
		// So, we check for the prefix.
		if !strings.Contains(rr.Body.String(), "Error sending shutdown command:") {
			t.Errorf("handler returned unexpected body: got '%v', want to contain 'Error sending shutdown command:'",
				rr.Body.String())
		}
		// Additionally, ensure it's not an empty error message part
		if strings.TrimSpace(strings.Replace(rr.Body.String(), "Error sending shutdown command:", "", 1)) == "" {
			t.Errorf("handler returned an empty error message part: got '%v'", rr.Body.String())
		}
	})
}
