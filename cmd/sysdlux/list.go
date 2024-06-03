package sysdlux

import (
	"flag"
	"fmt"
	"os"
)

const svcPath = "/Users/apb/workspaces/github.com/alexander-bergeron/go-systemdelux/svcs/"

func ListCommand(args []string) {
	fs := flag.NewFlagSet("list", flag.ExitOnError)
	// detached := fs.Bool("detached", false, "Run in detached mode")
	// detachedShort := fs.Bool("d", false, "Run in detached mode (short flag)")

	fs.Parse(args)

	files, err := os.ReadDir(svcPath)
	if err != nil {
		fmt.Errorf("Failed to read directory: %v", err)
	}

	for _, file := range files {
		fmt.Println(file)
	}

}
