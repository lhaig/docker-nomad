package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

type Metadata struct {
	SchemaVersion    string `json:"SchemaVersion"`
	Vendor           string `json:"Vendor"`
	Version          string `json:"Version"`
	ShortDescription string `json:"ShortDescription"`
	Experimental     bool   `json:"Experimental"`
	URL              string `json:"URL"`
}

func main() {
	if len(os.Args) > 1 && os.Args[1] == "docker-cli-plugin-metadata" {
		metadata := Metadata{
			SchemaVersion:    "0.1.0",
			Vendor:           "lhaig",
			Version:          "0.1.0",
			ShortDescription: "Docker CLI plugin to interact with Nomad",
			Experimental:     true,
			URL:              "https://github.com/lhaig/docker-nomad",
		}

		jsonOutput, err := json.MarshalIndent(metadata, "", "  ")
		if err != nil {
			fmt.Println("Error generating metadata:", err)
			os.Exit(1)
		}

		fmt.Println(string(jsonOutput))
		return
	}

	// Check if nomad binary is accessible
	path, err := exec.LookPath("nomad")
	if err != nil {
		fmt.Println("Nomad binary not found in PATH.")
		os.Exit(1)
	}

	fmt.Printf("Found Nomad binary at: %s\n", path)

	// Exclude the first argument (plugin name) and check if the second argument is "nomad"
	args := os.Args[1:]

	// If the first argument is "nomad", remove it from the arguments
	if len(args) > 0 && args[0] == "nomad" {
		args = args[1:]
	}

	// If no arguments are passed, default to "-help"
	if len(args) == 0 {
		fmt.Println("No arguments provided. Defaulting to -help.")
		args = append(args, "-help")
	}

	// Debug: Print the final arguments being passed to the nomad binary
	fmt.Printf("Arguments passed to Nomad: %v\n", args)

	// Prepare to run the nomad command with the provided arguments
	cmd := exec.Command("nomad", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Run the nomad command and check for errors
	err = cmd.Run()
	if err != nil {
		// Check if it's an exit status error
		if exitError, ok := err.(*exec.ExitError); ok {
			// Check the exit status
			if status, ok := exitError.Sys().(syscall.WaitStatus); ok {
				fmt.Printf("Nomad command failed with exit status: %d\n", status.ExitStatus())
			} else {
				fmt.Printf("Nomad command failed: %v\n", err)
			}
		} else {
			// General error (not related to exit status)
			fmt.Printf("Failed to execute nomad command: %v\n", err)
		}
		os.Exit(1)
	}
}
