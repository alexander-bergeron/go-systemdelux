package sysdlux

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/alexander-bergeron/go-systemdelux/monitor"
)

// Health check interval in minutes
const healthCheckInterval = 1

// const svcPath = "/Users/apb/workspaces/github.com/alexander-bergeron/go-systemdelux/svcs/"

func DaemonCommand(args []string) {
	fs := flag.NewFlagSet("daemon", flag.ExitOnError)
	detached := fs.Bool("detached", false, "Run in detached mode")
	detachedShort := fs.Bool("d", false, "Run in detached mode (short flag)")

	fs.Parse(args)

	if *detached || *detachedShort {
		// runDetached("sleep", args)
		fmt.Println("detached flag.")
		return
	} else {
		runOnce()
		// for {
		// 	fmt.Println("Running Health Check.")
		// 	runOnce()
		// 	time.Sleep(1 * time.Minute)
		// }
	}
}

func _main() {
	daemonFlag := flag.Bool("daemon", false, "Run as a daemon")
	flag.Parse()

	if *daemonFlag {
		// Daemonize the process
		runAsDaemon()
	} else {
		// Run health check once
		runOnce()
	}
}

func runAsDaemon() {
	// Fork the process
	cmd := exec.Command(os.Args[0], os.Args[1:]...)
	// cmd.Args = append(cmd.Args, "-daemon=true")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = os.Environ()

	if err := cmd.Start(); err != nil {
		log.Fatalf("Failed to start daemon process: %v", err)
	}

	// Detach from the parent process
	fmt.Printf("Daemon process started with PID %d\n", cmd.Process.Pid)
	os.Exit(0)
}

func runOnce() {
	var mon monitor.Monitor
	mon.LoadFromDirectory(svcPath)
	mon.CheckServices()
}

func runDaemon() {
	var mon monitor.Monitor
	mon.LoadFromDirectory(svcPath)

	ticker := time.NewTicker(time.Minute * healthCheckInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			mon.CheckServices()
		}
	}
}

func runDetached(command string, args []string) {
	cmd := exec.Command(os.Args[0], append([]string{command}, args...)...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		log.Fatalf("Failed to start process in detached mode: %v", err)
	}
	fmt.Printf("Running in detached mode with PID %d\n", cmd.Process.Pid)
	os.Exit(0)
}

// func init() {
// 	// Handle daemon flag explicitly if set
// 	if len(os.Args) > 1 && os.Args[1] == "-daemon=false" {
// 		runDaemon()
// 		os.Exit(0)
// 	}
// }
