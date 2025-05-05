package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/popovpsk/llama-manager/config"
	"github.com/popovpsk/llama-manager/processmanager"
	"github.com/popovpsk/llama-manager/server"
)

func main() {
	configPath := flag.String("config", "config.yaml", "Path to config file")
	flag.Parse()

	cfg, err := config.Load(*configPath)
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	pm := processmanager.NewProcessManager()
	// Pass the config path to the server
	srv := server.NewServer(cfg, pm, *configPath)

	// Handle graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-c
		log.Println("Shutting down server...")
		pm.StopCurrent()
		srv.Shutdown()
		os.Exit(0)
	}()

	log.Println("Starting server on :7000")
	if err := srv.Start(":7000"); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
