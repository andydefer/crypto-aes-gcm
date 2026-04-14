// Package ui provides terminal user interface utilities for aescryptool.
//
// This package handles all user interaction including:
//   - Colored output for different message types (info, success, error, warning)
//   - Progress bars for long-running operations
//   - Interactive prompts for file paths, passwords, and confirmations
//   - Banner displays for interactive mode
//
// All UI functions are designed to work consistently across different terminals
// and operating systems.
package ui

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"runtime"
	"strings"

	"github.com/andydefer/crypto-aes-gcm/internal/lang"
	"github.com/andydefer/crypto-aes-gcm/pkg/cryptolib"
	"github.com/manifoldco/promptui"
)

// PromptOperation displays a selection menu and returns the user's choice.
func PromptOperation() string {
	prompt := promptui.Select{
		Label: lang.T(lang.UIPromptOperationLabel),
		Items: []string{
			lang.T(lang.UIPromptEncryptOption),
			lang.T(lang.UIPromptDecryptOption),
			lang.T(lang.UIPromptExitOption),
		},
		Size: 5,
	}

	idx, _, err := prompt.Run()
	if err != nil {
		fmt.Println()
		SuccessColor.Println(lang.T(lang.UIPromptGoodbye))
		fmt.Println()
		os.Exit(0)
	}

	switch idx {
	case 0:
		return "encrypt"
	case 1:
		return "decrypt"
	default:
		return "exit"
	}
}

// PromptFilePath asks the user for a file path with optional validation.
func PromptFilePath(label string, mustExist bool, defaultValue string) string {
	for {
		prompt := promptui.Prompt{
			Label: label,
		}
		if defaultValue != "" {
			prompt.Default = defaultValue
		}

		result, err := prompt.Run()
		if err != nil {
			if err == promptui.ErrInterrupt {
				return ""
			}
			fmt.Println()
			SuccessColor.Println(lang.T(lang.UIPromptGoodbye))
			fmt.Println()
			os.Exit(0)
		}

		result = strings.TrimSpace(result)

		if result == "" && defaultValue != "" {
			result = strings.TrimSpace(defaultValue)
		}

		if result == "" {
			ErrorColor.Println(lang.T(lang.UIPromptPathEmpty))
			continue
		}

		if mustExist {
			if _, err := os.Stat(result); os.IsNotExist(err) {
				ErrorColor.Printf(lang.T(lang.UIPromptPathNotExist), result)
				fmt.Println()
				continue
			}
		}

		SuccessColor.Printf(lang.T(lang.UIPromptPathSuccess), result)
		return result
	}
}

// PromptPassword asks the user for a password with masked input.
func PromptPassword(label string, needValidation bool) string {
	for {
		prompt := promptui.Prompt{
			Label: label,
			Mask:  '*',
		}

		result, err := prompt.Run()
		if err != nil {
			if err == promptui.ErrInterrupt {
				return ""
			}
			fmt.Println()
			SuccessColor.Println(lang.T(lang.UIPromptGoodbye))
			fmt.Println()
			os.Exit(0)
		}

		result = strings.TrimSpace(result)

		if needValidation {
			if len(result) < 8 {
				ErrorColor.Println(lang.T(lang.UIPromptPasswordMinLength))
				continue
			}
			if !regexp.MustCompile(`[A-Z]`).MatchString(result) {
				ErrorColor.Println(lang.T(lang.UIPromptPasswordUppercase))
				continue
			}
			if !regexp.MustCompile(`[a-z]`).MatchString(result) {
				ErrorColor.Println(lang.T(lang.UIPromptPasswordLowercase))
				continue
			}
			if !regexp.MustCompile(`[0-9]`).MatchString(result) {
				ErrorColor.Println(lang.T(lang.UIPromptPasswordDigit))
				continue
			}
		}

		SuccessColor.Printf(lang.T(lang.UIPromptPasswordSuccess), strings.Repeat("*", len(result)))
		return result
	}
}

// PromptWorkers asks the user for the number of parallel encryption workers.
func PromptWorkers() int {
	maxWorkers := runtime.NumCPU() * 2

	for {
		prompt := promptui.Prompt{
			Label:   fmt.Sprintf(lang.T(lang.UIPromptWorkersLabel), cryptolib.DefaultWorkers(), maxWorkers),
			Default: fmt.Sprintf("%d", cryptolib.DefaultWorkers()),
		}

		result, err := prompt.Run()
		if err != nil {
			if err == promptui.ErrInterrupt {
				return cryptolib.DefaultWorkers()
			}
			fmt.Println()
			SuccessColor.Println(lang.T(lang.UIPromptGoodbye))
			fmt.Println()
			os.Exit(0)
		}

		result = strings.TrimSpace(result)

		if result == "" {
			SuccessColor.Printf(lang.T(lang.UIPromptWorkersSuccess), cryptolib.DefaultWorkers())
			return cryptolib.DefaultWorkers()
		}

		var w int
		_, err = fmt.Sscanf(result, "%d", &w)
		if err != nil || w < 1 {
			ErrorColor.Println(lang.T(lang.UIPromptWorkersInvalid))
			continue
		}
		if w > maxWorkers {
			ErrorColor.Printf(lang.T(lang.UIPromptWorkersMax), maxWorkers)
			continue
		}

		SuccessColor.Printf(lang.T(lang.UIPromptWorkersSuccess), w)
		return w
	}
}

// PromptConfirm asks the user for a yes/no confirmation.
func PromptConfirm(label string, defaultValue bool) bool {
	defaultDisplay := "Y/n"
	if !defaultValue {
		defaultDisplay = "y/N"
	}

	for {
		fmt.Printf(lang.T(lang.UIPromptConfirmLabel), label, defaultDisplay)

		reader := bufio.NewReader(os.Stdin)
		result, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println()
			SuccessColor.Println(lang.T(lang.UIPromptGoodbye))
			fmt.Println()
			os.Exit(0)
		}

		result = strings.TrimSpace(strings.ToLower(result))

		if result == "" {
			return defaultValue
		}

		if result == "y" || result == "yes" || result == "o" || result == "oui" {
			return true
		}
		if result == "n" || result == "no" || result == "non" {
			return false
		}

		ErrorColor.Println(lang.T(lang.UIPromptConfirmInvalid))
	}
}
