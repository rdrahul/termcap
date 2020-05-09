package utils

import (
	"fmt"
	"os"
)

//Er : prints the error before exiting
func Er(msg interface{}) {
	fmt.Println("Error:", msg)
	os.Exit(1)
}

// GetShell : return current active shell
func GetShell() string {

	if os.Getenv("SHELL") == "" {
		return "bin/bash"
	}
	return os.Getenv("SHELL")

}
