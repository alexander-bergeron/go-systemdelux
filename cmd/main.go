package main

import (
	"fmt"
	"os"

	"github.com/alexander-bergeron/go-systemdelux/cmd/sysdlux"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Expected 'sleep' or 'list' subcommands")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "list":
		sysdlux.ListCommand(os.Args[2:])
	case "add":
		sysdlux.AddCommand(os.Args[2:])
	case "daemon":
		sysdlux.DaemonCommand(os.Args[2:])
	default:
		fmt.Println("Expected 'list' subcommands")
		os.Exit(1)
	}
}
