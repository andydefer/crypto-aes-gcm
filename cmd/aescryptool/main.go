// Package main provides the entry point for aesaesaescryptool CLI.
//
// Aesaesaescryptool is a secure file encryption utility that uses AES-256-GCM with
// Argon2id key derivation and parallel streaming encryption.
//
// Usage:
//
//	aesaesaescryptool encrypt [input] [output] --pass <password>
//	aesaesaescryptool decrypt [input] [output] --pass <password>
//	aesaesaescryptool interact
//	aesaesaescryptool version
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
