package server

import (
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"os/exec"

	"github.com/popovpsk/llama-manager/config"
	"github.com/popovpsk/llama-manager/processmanager"
	"github.com/popovpsk/llama-manager/templates"
)

var execCommand = exec.Command

type Server struct {
	cfg        *config.Config
	pm         *processmanager.ProcessManager
	tmpl       *template.Template
	configPath string
}

func NewServer(cfg *config.Config, pm *processmanager.ProcessManager, configPath string) *Server {
	tmpl := template.Must(template.ParseFS(templates.TemplateFS, "index.html"))
	return &Server{
		cfg:        cfg,
		pm:         pm,
		tmpl:       tmpl,
		configPath: configPath,
	}
}

func (s *Server) Start(addr string) error {
	http.HandleFunc("/", s.handleIndex)
	http.HandleFunc("/run", s.handleRun)
	http.HandleFunc("/stop", s.handleStop)
	http.HandleFunc("/config", s.handleConfig)
	http.HandleFunc("/shutdown", s.handleShutdown)
	http.HandleFunc("/current-model", s.handleCurrentModel)
	return http.ListenAndServe(addr, nil)
}

func (s *Server) Shutdown() {
	s.pm.StopCurrent()
}

func (s *Server) handleIndex(w http.ResponseWriter, r *http.Request) {
	s.tmpl.Execute(w, s.cfg.Runs)
}

func (s *Server) handleRun(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	run := s.cfg.GetRun(name)
	if run == nil {
		http.Error(w, "Run not found", http.StatusNotFound)
		return
	}

	cmd := exec.Command("bash", "-c", run.Params.BuildCommand())
	err := s.pm.StartProcess(cmd)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error starting run: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Started run: %s", name)
}

func (s *Server) handleStop(w http.ResponseWriter, r *http.Request) {
	s.pm.StopCurrent()
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Stopped current run")
}

// handleConfig reads and returns the content of the configuration file.
func (s *Server) handleConfig(w http.ResponseWriter, r *http.Request) {
	data, err := os.ReadFile(s.configPath)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error reading config file: %v", err), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/yaml; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

// handleCurrentModel fetches model information from the LLaMA server and proxies it.
func (s *Server) handleCurrentModel(w http.ResponseWriter, r *http.Request) {
	// Make a request to the LLaMA server's /v1/models endpoint
	resp, err := http.Get("http://localhost:8080/v1/models")
	if err != nil {
		http.Error(w, fmt.Sprintf("Error contacting LLaMA server: %v", err), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Copy headers from the LLaMA server's response to our response
	for key, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}

	// Set the status code from the LLaMA server's response
	w.WriteHeader(resp.StatusCode)

	// Copy the body from the LLaMA server's response to our response
	if _, err := io.Copy(w, resp.Body); err != nil {
		// Log the error on the server side, client might have already received headers/status
		fmt.Printf("Error copying LLaMA server response body: %v\n", err)
		// Avoid writing an http.Error here if headers/status already sent
	}
}

// handleShutdown attempts to shut down the PC.
func (s *Server) handleShutdown(w http.ResponseWriter, r *http.Request) {
	// For Ubuntu/Linux, 'systemctl poweroff' or 'shutdown -h now' are common commands.
	// For macOS, 'osascript -e 'tell app "System Events" to shut down''.
	// For Windows, 'shutdown /s /t 0'.
	// Using 'systemctl poweroff' for Ubuntu.
	// IMPORTANT: Ensure the server process has permissions to execute this (e.g., via sudoers).
	// This is a potentially dangerous operation.
	cmd := execCommand("systemctl", "poweroff")
	err := cmd.Start()                          // Use Start for non-blocking, or Run for blocking
	if err != nil {
		http.Error(w, fmt.Sprintf("Error sending shutdown command: %v", err), http.StatusInternalServerError)
		return
	}

	// It's better not to wait for the command to complete with cmd.Wait()
	// because the server itself might be terminated as part of the shutdown.
	// We're just initiating the shutdown.

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Shutdown command issued to PC. The system should shut down shortly.")
}
