package monitor

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/alexander-bergeron/go-systemdelux/internal/go-ps"
	"gopkg.in/yaml.v2"
)

type Service struct {
	Name        string            `yaml:"name" json:"name"`
	Restart     bool              `yaml:"restart" json:"restart"`
	Owner       string            `yaml:"owner" json:"owner"`
	Hostname    string            `yaml:"hostname" json:"hostname"`
	Port        int               `yaml:"port" json:"port"`
	Command     []string          `yaml:"command" json:"command"`
	Environment map[string]string `yaml:"environment" json:"environment"`
	Flags       map[string]string `yaml:"flags" json:"flags"`
	Startup     string            `yaml:"startup" json:"startup"`
	Shutdown    string            `yaml:"shutdown" json:"shutdown"`
	Healthcheck string            `yaml:"healthcheck" json:"healthcheck"`
	Log         string            `yaml:"log" json:"log"`
}

type Monitor struct {
	Services []Service
}

func (m *Monitor) CheckServices() {

	for _, svc := range m.Services {
		svc.CheckHealth()
	}

}

func (m *Monitor) LoadFromDirectory(dirPath string) {

	files, err := os.ReadDir(dirPath)
	if err != nil {
		fmt.Printf("Failed to read directory: %v", err)
	}

	for _, file := range files {
		if !file.IsDir() {
			filePath := filepath.Join(dirPath, file.Name())

			data, err := os.ReadFile(filePath)
			if err != nil {
				fmt.Printf("Failed to read file %s: %v", filePath, err)
				continue
			}

			var svc Service
			err = yaml.Unmarshal(data, &svc)
			if err != nil {
				fmt.Printf("Failed to unmarshal YAML file %s: %v", filePath, err)
				continue
			}

			m.Services = append(m.Services, svc)
			fmt.Printf("Loaded service from %s: %+v\n", filePath, svc)
		}
	}

}

func (s *Service) CheckHealth() (bool, error) {

	if !s.Restart {
		return false, nil
	}

	running, err := s.isRunning()

	if err != nil {
		return false, nil
	}

	isOwner, err := s.isOwner()
	if err != nil {
		return false, nil
	}

	onHost, err := s.onHost()
	if err != nil {
		return false, nil
	}

	if !running && isOwner && onHost {
		fmt.Println("Restarting Service.")
		err := s.restart()
		return true, err
	}

	return false, nil

}

func (s *Service) restart() error {

	if len(s.Command) == 0 {
		return fmt.Errorf("command is empty")
	}

	err := s.setVars()
	if err != nil {
		return err
	}

	cmd := exec.Command(s.Command[0], s.Command[1:]...)

	// Disown the process
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
	}

	// Set the process to write its output to a log file
	logFile, err := os.OpenFile(s.Log, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open log file: %v", err)
	}
	defer logFile.Close()

	cmd.Stdout = logFile
	cmd.Stderr = logFile

	// Start the command
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start command: %v", err)
	}

	fmt.Printf("Started command with PID %d\n", cmd.Process.Pid)
	return nil

}

func (s *Service) isRunning() (bool, error) {

	processes, err := ps.GetProcesses()
	if err != nil {
		fmt.Println("Error:", err)
		return false, err
	}

	for _, process := range processes {
		cmd := strings.Join(s.Command, " ")
		if cmd == process.Command {
			fmt.Println("Process is running.")
			return true, nil
		}
	}

	return false, nil

}

func (s *Service) isOwner() (bool, error) {

	processes, err := ps.GetProcesses()
	if err != nil {
		fmt.Println("Error:", err)
		return false, err
	}

	for _, process := range processes {
		if s.Owner == process.User {
			fmt.Println("Current User owns the process.")
			return true, nil
		}
	}

	return false, nil
}

func (s *Service) onHost() (bool, error) {

	currentHost, err := os.Hostname()
	if err != nil {
		fmt.Println("Error: ", err)
		return false, err
	}

	// fmt.Println(s.Hostname)
	// fmt.Println(currentHost)

	if s.Hostname == currentHost {
		fmt.Println("Current host is the correct host.")
		return true, nil
	}

	return false, nil
}

func (s *Service) setVars() error {
	for key, value := range s.Environment {
		err := os.Setenv(key, value)
		if err != nil {
			return fmt.Errorf("failed to set environment variable %s: %w", key, err)
		}
	}
	return nil
}

func all(conditions ...bool) bool {
	for _, condition := range conditions {
		if !condition {
			return false
		}
	}
	return true
}
