// Package main provides the entry point for cryptool CLI.
package main

import (
	"os"

	"github.com/andydefer/crypto-aes-gcm/internal/cli"
)

func main() {
	if err := cli.Execute(); err != nil {
		os.Exit(1)
	}
}
