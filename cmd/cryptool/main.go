// Package main provides the entry point for cryptool CLI.
//
// Cryptool is a secure file encryption utility that uses AES-256-GCM with
// Argon2id key derivation and parallel streaming encryption.
//
// Usage:
//
//	cryptool encrypt [input] [output] --pass <password>
//	cryptool decrypt [input] [output] --pass <password>
//	cryptool interact
//	cryptool version
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
