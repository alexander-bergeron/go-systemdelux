package ps

import (
	"bufio"
	"os/exec"
	"strconv"
	"strings"
)

type Process struct {
	User    string
	PID     int
	CPU     float64
	MEM     float64
	VSZ     int
	RSS     int
	TT      string
	STAT    string
	Started string
	Time    string
	Command string
}

func (p *Process) Kill() error {
	cmd := exec.Command("kill", "-9", strconv.Itoa(p.PID))

	err := cmd.Run()
	return err
}

func GetProcesses() ([]Process, error) {
	cmd := exec.Command("ps", "aux")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(strings.NewReader(string(output)))
	var processes []Process

	// Skip the header line
	scanner.Scan()

	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)
		if len(fields) < 11 {
			continue // skip malformed lines
		}

		pid, _ := strconv.Atoi(fields[1])
		cpu, _ := strconv.ParseFloat(fields[2], 64)
		mem, _ := strconv.ParseFloat(fields[3], 64)
		vsz, _ := strconv.Atoi(fields[4])
		rss, _ := strconv.Atoi(fields[5])

		process := Process{
			User:    fields[0],
			PID:     pid,
			CPU:     cpu,
			MEM:     mem,
			VSZ:     vsz,
			RSS:     rss,
			TT:      fields[6],
			STAT:    fields[7],
			Started: fields[8],
			Time:    fields[9],
			Command: strings.Join(fields[10:], " "), // Command can have spaces
		}

		processes = append(processes, process)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return processes, nil
}

// Function to filter the list of processes by user
func FilterByUser(processes []Process, user string) []Process {

	var filtered []Process
	for _, process := range processes {
		if process.User == user {
			filtered = append(filtered, process)
			// break
		}
	}
	return filtered
}

// Function to filter the list of processes by command
func FilterByCommand(processes []Process, command string) []Process {

	var filtered []Process
	for _, process := range processes {
		if strings.Contains(process.Command, command) {
			filtered = append(filtered, process)
			// break
		}
	}

	return filtered
}
