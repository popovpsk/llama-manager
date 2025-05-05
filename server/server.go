package server

import (
	"fmt"
	"html/template"
	"net/http"
	"os" // Added for reading the config file
	"os/exec"
	"github.com/popovpsk/llama-manager/config"
	"github.com/popovpsk/llama-manager/processmanager"
	"github.com/popovpsk/llama-manager/templates"
)

type Server struct {
	cfg        *config.Config
	pm         *processmanager.ProcessManager
	tmpl       *template.Template
	configPath string // Added config path
}

func NewServer(cfg *config.Config, pm *processmanager.ProcessManager, configPath string) *Server {
	tmpl := template.Must(template.ParseFS(templates.TemplateFS, "index.html"))
	return &Server{
		cfg:        cfg,
		pm:         pm,
		tmpl:       tmpl,
		configPath: configPath, // Store config path
	}
}

func (s *Server) Start(addr string) error {
	http.HandleFunc("/", s.handleIndex)
	http.HandleFunc("/run", s.handleRun)
	http.HandleFunc("/stop", s.handleStop)
	http.HandleFunc("/config", s.handleConfig) // Added config handler route
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

	cmd := exec.Command("bash", "-c", run.Cmd)
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
	w.Header().Set("Content-Type", "text/yaml; charset=utf-8") // Indicate YAML content type
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}
