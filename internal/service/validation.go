package service

import (
	"fmt"
	"os"
	"runtime"

	"github.com/andydefer/crypto-aes-gcm/internal/ui"
	"github.com/andydefer/crypto-aes-gcm/pkg/cryptolib"
	"github.com/manifoldco/promptui"
)

// ValidateWorkerCount ensures the worker count is within reasonable bounds.
func ValidateWorkerCount(requested int, quiet bool) int {
	if requested <= 0 {
		return cryptolib.DefaultWorkers
	}
	maxWorkers := runtime.NumCPU() * 2
	if requested > maxWorkers {
		if !quiet {
			ui.WarningColor.Printf("⚠️ Workers réduit à %d\n", maxWorkers)
		}
		return maxWorkers
	}
	return requested
}

// ValidateInputFile checks if the input file exists and is accessible.
func ValidateInputFile(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return fmt.Errorf("fichier '%s' inexistant", path)
	}
	return nil
}

// CheckFileExists checks if a file exists.
func CheckFileExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

// CheckOverwrite prompts for confirmation when output file exists (non-interactive).
func CheckOverwrite(output string, force bool) error {
	if force {
		return nil
	}
	if _, err := os.Stat(output); err == nil {
		prompt := promptui.Prompt{
			Label:     fmt.Sprintf("Fichier '%s' existe. Écraser ?", output),
			IsConfirm: true,
			Default:   "n",
		}
		result, err := prompt.Run()
		if err != nil || (result != "y" && result != "Y") {
			return fmt.Errorf("annulé")
		}
	}
	return nil
}
