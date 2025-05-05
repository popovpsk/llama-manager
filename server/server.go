package server

import (
	"fmt"
	"html/template"
	"net/http"
	"os/exec"
	"path/filepath"

	"github.com/popovpsk/llama-manager/config"
	"github.com/popovpsk/llama-manager/processmanager"
)

type Server struct {
	cfg  *config.Config
	pm   *processmanager.ProcessManager
	tmpl *template.Template
}

func NewServer(cfg *config.Config, pm *processmanager.ProcessManager) *Server {
	tmpl := template.Must(template.ParseFiles(filepath.Join("templates", "index.html")))
	return &Server{
		cfg:  cfg,
		pm:   pm,
		tmpl: tmpl,
	}
}

func (s *Server) Start(addr string) error {
	http.HandleFunc("/", s.handleIndex)
	http.HandleFunc("/run", s.handleRun)
	http.HandleFunc("/stop", s.handleStop)
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
